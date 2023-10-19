package service

import (
	"context"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
)

type AddAccountInput struct {
	validationErrors map[string]string
	Email            string
	Name             string
	Avatar           string
}

type UpdateAccountInput struct {
	Id     string
	Email  string
	Name   string
	Avatar string
}

// Returns an account by searching with provided id
//
// Errors returned: ErrFailedValidation, ErrNoRows
func (s *Service) GetAccountById(ctx context.Context, id string) (*types.Account, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrFailedValidation
	}

	var acc types.Account
	query := "SELECT id, email, name, avatar, deleted, created_at, updated_at FROM accounts WHERE id::text=$1"
	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&acc.Id, &acc.Email, &acc.Name, &acc.Avatar, &acc.Deleted, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

// Returns a slice accounts by searching with provided id
func (s *Service) GetAllAccounts(ctx context.Context) ([]types.Account, error) {
	accs := make([]types.Account, 0)

	rows, err := s.DB.QueryContext(ctx, "SELECT id, email, name, avatar, deleted, created_at, updated_at FROM accounts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var acc types.Account
		err := rows.Scan(&acc.Id, &acc.Email, &acc.Name, &acc.Avatar, &acc.Deleted, &acc.CreatedAt, &acc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if !acc.Deleted {
			accs = append(accs, acc)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accs, nil
}

// Inserts a new record in db
//
// Errors returned: ErrFailedValidation, ErrFailedToUpdate
func (s *Service) AddAccount(ctx context.Context, input *AddAccountInput) error {
	input.validationErrors = make(map[string]string)

	if input.Email == "" {
		input.validationErrors["email"] = "must be provided"
	}

	if input.Name == "" {
		input.validationErrors["name"] = "must be provided"
	}

	if len(input.validationErrors) > 0 {
		return ErrFailedValidation
	}

	res, err := s.DB.ExecContext(ctx, "INSERT INTO accounts (name, email, avatar) VALUES ($1, $2, $3)", input.Name, input.Email, input.Avatar)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}

// Updates an account with provided data
//
// Errors returned: ErrFailedValidation, ErrFailedToUpdate
func (s *Service) UpdateAccount(ctx context.Context, input *UpdateAccountInput) error {
	if _, err := uuid.Parse(input.Id); err != nil {
		return ErrFailedValidation
	}

	query := "UPDATE accounts SET name=COALESCE(NULLIF($1, ''), name), email=COALESCE(NULLIF($2, ''), email), avatar=COALESCE(NULLIF($3, ''), avatar) WHERE id=$4"
	res, err := s.DB.ExecContext(ctx, query, input.Name, input.Email, input.Avatar, input.Id)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}

// Delete an account with provided id
//
// Errors returned: ErrFailedValidation, ErrFailedToUpdate
func (s *Service) DeleteAccount(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrFailedValidation
	}
	res, err := s.DB.ExecContext(ctx, "UPDATE accounts SET deleted=true WHERE id=$1", id)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra < 1 {
		return ErrFailedToUpdate
	}

	return nil
}
