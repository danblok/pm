package service

import (
	"context"
	"testing"

	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestGetProjectById(t *testing.T) {
	s, cleanup := setupService(t)

	pId := uuid.NewString()
	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantErr error
		want    *types.Project
		input   string
	}{
		"non-existent": {
			input:   uuid.NewString(),
			wantErr: ErrNotFound,
			want:    nil,
		},
		"invalid id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
			want:    nil,
		},
		"existing": {
			input:   pId,
			wantErr: nil,
			want: &types.Project{
				Id:      pId,
				Name:    "Existing project",
				OwnerId: owner.Id,
			},
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", tt.want.Id, tt.want.Name, tt.want.Description, tt.want.OwnerId)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			got, err := s.GetProjectById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetProjectById() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Project{}, "CreatedAt", "UpdatedAt", "Owner")); diff != "" {
				t.Fatalf("GetProjectById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetProjectsByOwnerId(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantErr error
		input   string
		want    []types.Project
	}{
		"invalid owner id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
			want:    []types.Project{},
		},
		"2 projects": {
			input:   owner.Id,
			wantErr: nil,
			want: []types.Project{
				{
					Id:      uuid.NewString(),
					Name:    "Project 1",
					OwnerId: owner.Id,
				},
				{
					Id:      uuid.NewString(),
					Name:    "Project 1",
					OwnerId: owner.Id,
				},
			},
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			for _, p := range tt.want {
				_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
				if err != nil {
					t.Fatal(ErrFailedToPrepareTest, err)
				}
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			got, err := s.GetProjectsByOwnerId(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetProjectsByOwnerId() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Project{}, "CreatedAt", "UpdatedAt", "Owner")); diff != "" {
				t.Fatalf("GetProjectsByOwnerId() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAddProject(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		input   *AddProjectInput
		wantErr error
	}{
		"invalid owner id": {
			input: &AddProjectInput{
				Name:    "project",
				OwnerId: "invalid-id",
			},
			wantErr: ErrFailedValidation,
		},
		"invalid name": {
			input: &AddProjectInput{
				OwnerId: owner.Id,
			},
			wantErr: ErrFailedValidation,
		},
		"non-existent owner id": {
			input: &AddProjectInput{
				Name:    "project",
				OwnerId: uuid.NewString(),
			},
			wantErr: ErrInternal,
		},
		"sucsessfull add": {
			input: &AddProjectInput{
				Name:    "project",
				OwnerId: owner.Id,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			err := s.AddProject(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("AddProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	p := types.Project{
		Id:      uuid.NewString(),
		Name:    "Project",
		OwnerId: owner.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   *UpdateProjectInput
	}{
		"invalid id": {
			input: &UpdateProjectInput{
				Id:   "invalid-id",
				Name: "New project name",
			},
			wantErr: ErrFailedValidation,
		},
		"non-existent project id": {
			input: &UpdateProjectInput{
				Id:   uuid.NewString(),
				Name: "New project name",
			},
			wantErr: ErrFailedToUpdate,
		},
		"sucsessfull update": {
			input: &UpdateProjectInput{
				Id:   p.Id,
				Name: "New project name",
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			err := s.UpdateProject(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("UpdateProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteProjectById(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	p := types.Project{
		Id:      uuid.NewString(),
		Name:    "Project",
		OwnerId: owner.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   string
	}{
		"invalid id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
		},
		"non-existent project id": {
			input:   uuid.NewString(),
			wantErr: ErrFailedToUpdate,
		},
		"sucsessfull delete": {
			input:   p.Id,
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			ctx := context.Background()
			err := s.DeleteProjectById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("UpdateProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
