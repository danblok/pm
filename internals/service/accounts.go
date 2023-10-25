package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
)

type AddAccountInput struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type UpdateAccountInput struct {
	Id     string `param:"id"`
	Email  string `json:"email,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// Errors returned: ErrFailedValidation, ErrInternal, ErrNotFound
func (s *Service) GetAccountById(ctx context.Context, id string) (*types.Account, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrFailedValidation
	}

	var acc types.Account
	query := "SELECT id, email, name, avatar, deleted, created_at, updated_at FROM accounts WHERE id::text=$1"
	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&acc.Id, &acc.Email, &acc.Name, &acc.Avatar, &acc.Deleted, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &acc, nil
}

// Returne errors: ErrInternal
func (s *Service) GetAllAccounts(ctx context.Context) ([]types.Account, error) {
	accs := make([]types.Account, 0)

	rows, err := s.DB.QueryContext(ctx, "SELECT id, email, name, avatar, deleted, created_at, updated_at FROM accounts")
	if err != nil {
		return nil, ErrInternal
	}

	for rows.Next() {
		var acc types.Account
		err := rows.Scan(&acc.Id, &acc.Email, &acc.Name, &acc.Avatar, &acc.Deleted, &acc.CreatedAt, &acc.UpdatedAt)
		if err != nil {
			return nil, ErrInternal
		}

		if !acc.Deleted {
			accs = append(accs, acc)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return accs, nil
}

// Errors returned: ErrFailedValidation, ErrFailedToUpdate, ErrFailedToInsert
func (s *Service) AddAccount(ctx context.Context, input *AddAccountInput) error {
	if input.Email == "" {
		return ErrFailedValidation
	}

	if input.Name == "" {
		return ErrFailedValidation
	}

	res, err := s.DB.ExecContext(ctx, "INSERT INTO accounts (name, email, avatar) VALUES ($1, $2, $3)", input.Name, input.Email, input.Avatar)
	if err != nil {
		return ErrInternal
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return ErrFailedToInsert
	}
	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}

// Errors returned: ErrFailedValidation, ErrFailedToUpdate, ErrInternal
func (s *Service) UpdateAccount(ctx context.Context, input *UpdateAccountInput) error {
	if _, err := uuid.Parse(input.Id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE accounts SET name=COALESCE(NULLIF($1, ''), name), email=COALESCE(NULLIF($2, ''), email), avatar=COALESCE(NULLIF($3, ''), avatar) WHERE id=$4"
	res, err := s.DB.ExecContext(ctx, query, input.Name, input.Email, input.Avatar, input.Id)
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

// Errors returned: ErrFailedValidation, ErrFailedToUpdate, ErrInternal
func (s *Service) DeleteAccountById(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrFailedValidation
	}
	res, err := s.DB.ExecContext(ctx, "UPDATE accounts SET deleted=true WHERE id=$1", id)
	if err != nil {
		return ErrInternal
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return ErrFailedToUpdate
	}
	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}
