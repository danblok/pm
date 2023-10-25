package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
)

type AddProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerId     string `json:"owner_id"`
}

type UpdateProjectInput struct {
	Id          string `param:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Returned errors: ErrFailedValidation, ErrInternal, ErrNotFound
func (s *Service) GetProjectById(ctx context.Context, id string) (*types.Project, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrFailedValidation
	}

	var pj types.Project
	acc := new(types.Account)
	pj.Owner = acc
	query := "SELECT * FROM projects WHERE id=$1 AND deleted=false"
	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&pj.Id, &pj.Name, &pj.Description, &pj.OwnerId, &pj.Deleted, &pj.CreatedAt, &pj.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &pj, nil
}

// Returned errors: ErrFailedValidation, ErrInternal
func (s *Service) GetProjectsByOwnerId(ctx context.Context, ownerId string) ([]types.Project, error) {
	pjs := make([]types.Project, 0)
	if _, err := uuid.Parse(ownerId); err != nil {
		return pjs, ErrFailedValidation
	}
	query := "SELECT * FROM projects WHERE owner_id=$1 AND deleted=false"
	rows, err := s.DB.QueryContext(ctx, query, ownerId)
	if err != nil {
		return nil, ErrInternal
	}

	for rows.Next() {
		var pj types.Project

		err = rows.Scan(&pj.Id, &pj.Name, &pj.Description, &pj.OwnerId, &pj.Deleted, &pj.CreatedAt, &pj.UpdatedAt)
		if err != nil {
			return nil, ErrInternal
		}

		pjs = append(pjs, pj)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return pjs, nil
}

// Returned errors: ErrFailedValidation, ErrInternal, ErrFailedToInsert
func (s *Service) AddProject(ctx context.Context, input *AddProjectInput) error {
	if input.Name == "" {
		return ErrFailedValidation
	}
	if _, err := uuid.Parse(input.OwnerId); err != nil {
		return ErrFailedValidation
	}

	query := "INSERT INTO projects (name, description, owner_id) VALUES ($1, $2, $3)"
	res, err := s.DB.ExecContext(ctx, query, input.Name, input.Description, input.OwnerId)
	if err != nil {
		return ErrInternal
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return ErrInternal
	}
	if ra < 1 {
		return ErrFailedToInsert
	}

	return nil
}

// Returned errors: ErrFailedValidation, ErrInternal, ErrFailedToUpdate
func (s *Service) UpdateProject(ctx context.Context, input *UpdateProjectInput) error {
	if _, err := uuid.Parse(input.Id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE projects SET name=COALESCE(NULLIF($1, ''), name), description=COALESCE(NULLIF($2, ''), description) WHERE id::text=$3"
	res, err := s.DB.ExecContext(ctx, query, input.Name, input.Description, input.Id)
	if err != nil {
		return ErrInternal
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return ErrInternal
	}

	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}

// Returned errors: ErrFailedValidation, ErrInternal, ErrFailedToUpdate
func (s *Service) DeleteProjectById(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE projects SET deleted=true WHERE id=$1"
	res, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return ErrInternal
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return ErrInternal
	}

	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}
