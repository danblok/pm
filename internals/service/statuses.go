package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
)

type AddStatusInput struct {
	Name      string `json:"name"`
	ProjectId string `json:"project_id"`
}

type UpdateStatusInput struct {
	Id   string `param:"id"`
	Name string `json:"name,omitempty"`
}

// Errors returned: ErrFailedValidation, ErrInternal, ErrNotFound
func (s *Service) GetStatusById(ctx context.Context, id string) (*types.Status, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrFailedValidation
	}

	var st types.Status
	query := "SELECT * FROM statuses WHERE id=$1 AND deleted=false"
	row := s.DB.QueryRow(query, id)
	err := row.Scan(&st.Id, &st.Name, &st.ProjectId, &st.Deleted, &st.CreatedAt, &st.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &st, nil
}

// Returned errors: ErrFailedValidation, ErrInternal
func (s *Service) GetStatusesByProjectId(ctx context.Context, pId string) ([]types.Status, error) {
	sts := make([]types.Status, 0)
	if _, err := uuid.Parse(pId); err != nil {
		return sts, ErrFailedValidation
	}
	query := "SELECT * FROM statuses WHERE project_id=$1 AND deleted=false"
	rows, err := s.DB.QueryContext(ctx, query, pId)
	if err != nil {
		return nil, ErrInternal
	}

	for rows.Next() {
		var st types.Status

		err = rows.Scan(&st.Id, &st.Name, &st.ProjectId, &st.Deleted, &st.CreatedAt, &st.UpdatedAt)
		if err != nil {
			return nil, ErrInternal
		}

		sts = append(sts, st)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return sts, nil
}

// Returned errors: ErrFailedValidation, ErrInternal, ErrFailedToInsert
func (s *Service) AddStatus(ctx context.Context, input *AddStatusInput) error {
	if input.Name == "" {
		return ErrFailedValidation
	}
	if _, err := uuid.Parse(input.ProjectId); err != nil {
		return ErrFailedValidation
	}

	query := "INSERT INTO statuses (name, project_id) VALUES ($1, $2)"
	res, err := s.DB.ExecContext(ctx, query, input.Name, input.ProjectId)
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
func (s *Service) UpdateStatus(ctx context.Context, input *UpdateStatusInput) error {
	if _, err := uuid.Parse(input.Id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE statuses SET name=COALESCE(NULLIF($1, ''), name) WHERE id::text=$2"
	res, err := s.DB.ExecContext(ctx, query, input.Name, input.Id)
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
func (s *Service) DeleteStatusById(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE statuses SET deleted=true WHERE id=$1"
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
