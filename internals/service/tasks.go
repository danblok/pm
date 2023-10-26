package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
)

type AddTaskInput struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	Name      string `json:"name"`
	ProjectId string `json:"project_id"`
	StatusId  string `json:"status_id"`
}

type UpdateTaskInput struct {
	Start    string `json:"start,omitempty"`
	End      string `json:"end,omitempty"`
	Id       string `param:"id"`
	Name     string `json:"name,omitempty"`
	StatusId string `json:"status_id"`
}

// Errors returned: ErrFailedValidation, ErrInternal, ErrNotFound
func (s *Service) GetTaskById(ctx context.Context, id string) (*types.Task, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrFailedValidation
	}

	var t types.Task
	query := "SELECT * FROM tasks WHERE id=$1 AND deleted=false"
	row := s.DB.QueryRow(query, id)
	err := row.Scan(&t.Id, &t.Name, &t.Start, &t.End, &t.StatusId, &t.ProjectId, &t.Deleted, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &t, nil
}

// Returned errors: ErrFailedValidation, ErrInternal
func (s *Service) GetTasksByProjectId(ctx context.Context, pId string) ([]types.Task, error) {
	ts := make([]types.Task, 0)
	if _, err := uuid.Parse(pId); err != nil {
		return ts, ErrFailedValidation
	}
	query := "SELECT * FROM tasks WHERE project_id=$1 AND deleted=false"
	rows, err := s.DB.QueryContext(ctx, query, pId)
	if err != nil {
		return nil, ErrInternal
	}

	for rows.Next() {
		var st types.Task

		err = rows.Scan(&st.Id, &st.Name, &st.Start, &st.End, &st.StatusId, &st.ProjectId, &st.Deleted, &st.CreatedAt, &st.UpdatedAt)
		if err != nil {
			return nil, ErrInternal
		}

		ts = append(ts, st)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return ts, nil
}

// Returned errors: ErrFailedValidation, ErrInternal
func (s *Service) GetTasksOfProjectByStatusId(ctx context.Context, pId, sId string) ([]types.Task, error) {
	ts := make([]types.Task, 0)
	if _, err := uuid.Parse(pId); err != nil {
		return ts, ErrFailedValidation
	}
	if _, err := uuid.Parse(sId); err != nil {
		return ts, ErrFailedValidation
	}
	query := "SELECT * FROM tasks WHERE project_id=$1 AND status_id=$2 AND deleted=false"
	rows, err := s.DB.QueryContext(ctx, query, pId, sId)
	if err != nil {
		return nil, ErrInternal
	}

	for rows.Next() {
		var st types.Task

		err = rows.Scan(&st.Id, &st.Name, &st.Start, &st.End, &st.StatusId, &st.ProjectId, &st.Deleted, &st.CreatedAt, &st.UpdatedAt)
		if err != nil {
			return nil, ErrInternal
		}

		ts = append(ts, st)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return ts, nil
}

// Returned errors: ErrFailedValidation, ErrInternal, ErrFailedToInsert
func (s *Service) AddTask(ctx context.Context, input *AddTaskInput) error {
	if input.Name == "" {
		return ErrFailedValidation
	}
	if _, err := uuid.Parse(input.ProjectId); err != nil {
		return ErrFailedValidation
	}
	if _, err := uuid.Parse(input.StatusId); err != nil {
		return ErrFailedValidation
	}
	start, err := time.Parse(time.DateTime, input.Start)
	if err != nil {
		return ErrFailedValidation
	}
	end, err := time.Parse(time.DateTime, input.End)
	if err != nil {
		return ErrFailedValidation
	}
	if end.Before(start) {
		return ErrFailedValidation
	}

	query := "INSERT INTO tasks (name, \"start\", \"end\", project_id, status_id) VALUES ($1, $2, $3, $4, $5)"
	res, err := s.DB.ExecContext(ctx, query, input.Name, start.UTC(), end.UTC(), input.ProjectId, input.StatusId)
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
func (s *Service) UpdateTask(ctx context.Context, input *UpdateTaskInput) error {
	if _, err := uuid.Parse(input.Id); err != nil {
		return ErrFailedValidation
	}
	if _, err := uuid.Parse(input.StatusId); err != nil {
		return ErrFailedValidation
	}
	start, err := time.Parse(time.DateTime, input.Start)
	if err != nil {
		return ErrFailedValidation
	}
	end, err := time.Parse(time.DateTime, input.End)
	if err != nil || end.Before(start) {
		return ErrFailedValidation
	}

	query := "UPDATE tasks SET name=COALESCE(NULLIF($1, ''), name), \"start\"=$2, \"end\"=$3, status_id=$4 WHERE id::text=$5"
	res, err := s.DB.ExecContext(ctx, query, input.Name, start, end, input.StatusId, input.Id)
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
func (s *Service) DeleteTaskById(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE tasks SET deleted=true WHERE id=$1"
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
