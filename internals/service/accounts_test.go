package service

import (
	"context"
	"testing"

	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestGetAccountById(t *testing.T) {
	s, cleanup := setupService(t)

	accId := uuid.NewString()
	tests := map[string]struct {
		wantErr error
		want    *types.Account
		input   string
	}{
		"existent": {
			input:   accId,
			wantErr: nil,
			want: &types.Account{
				Id:    accId,
				Name:  "username",
				Email: "username@test.com",
			},
		},
		"non-existent": {
			input:   uuid.NewString(),
			wantErr: ErrNotFound,
			want:    nil,
		},
		"invalid id": {
			input:   "invald-id",
			wantErr: ErrFailedValidation,
			want:    nil,
		},
	}

	for name, tt := range tests {
		if tt.want != nil {
			_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", tt.want.Id, tt.want.Email, tt.want.Name)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			got, err := s.GetAccountById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetAccountById() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Account{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("GetAccountById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetAllAccounts(t *testing.T) {
	s, cleanup := setupService(t)

	tests := map[string]struct {
		wantErr error
		want    []types.Account
	}{
		"2 projects": {
			wantErr: nil,
			want: []types.Account{
				{
					Id:    uuid.NewString(),
					Name:  "username 1",
					Email: "username1@test.com",
				},
				{
					Id:    uuid.NewString(),
					Name:  "username2",
					Email: "username2@test.com",
				},
			},
		},
		"0 projects": {
			wantErr: nil,
			want:    make([]types.Account, 0),
		},
	}

	for name, tt := range tests {
		for _, p := range tt.want {
			_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", p.Id, p.Email, p.Name)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			got, err := s.GetAllAccounts(ctx)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetAccountById() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Account{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("GetAccountById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAddAccount(t *testing.T) {
	s, cleanup := setupService(t)

	tests := map[string]struct {
		wantErr error
		input   *AddAccountInput
	}{
		"succsessfull add": {
			input: &AddAccountInput{
				Name:  "username",
				Email: "username@test.com",
			},
			wantErr: nil,
		},
		"invalid name": {
			input: &AddAccountInput{
				Name:  "",
				Email: "username@test.com",
			},
			wantErr: ErrFailedValidation,
		},
		"invalid email": {
			input: &AddAccountInput{
				Name:  "Project",
				Email: "",
			},
			wantErr: ErrFailedValidation,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			err := s.AddAccount(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("AddAccount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	s, cleanup := setupService(t)

	acc := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantErr error
		input   *UpdateAccountInput
	}{
		"succsessfull update": {
			input: &UpdateAccountInput{
				Id:    acc.Id,
				Name:  "New project",
				Email: "newusername@test.com",
			},
			wantErr: nil,
		},
		"non-existent id": {
			input: &UpdateAccountInput{
				Id:    uuid.NewString(),
				Name:  "New project",
				Email: "newusername@test.com",
			},
			wantErr: ErrFailedToUpdate,
		},
		"invalid id": {
			input: &UpdateAccountInput{
				Id:    "invalid-id",
				Name:  "",
				Email: "username@test.com",
			},
			wantErr: ErrFailedValidation,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			err := s.UpdateAccount(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("UpdateAccount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteAccountById(t *testing.T) {
	s, cleanup := setupService(t)

	acc := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantErr error
		input   string
	}{
		"succsessfull delete": {
			input:   acc.Id,
			wantErr: nil,
		},
		"non-existent": {
			input:   uuid.NewString(),
			wantErr: ErrFailedToUpdate,
		},
		"invalid id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			err := s.DeleteAccountById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("DeleteAccountById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
