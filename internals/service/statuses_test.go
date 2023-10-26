package service

import (
	"context"
	"testing"

	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestGetStatusById(t *testing.T) {
	s, cleanup := setupService(t)

	sId := uuid.NewString()
	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	project := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "Discription of the project",
		OwnerId:     owner.Id,
	}
	tests := map[string]struct {
		wantErr error
		want    *types.Status
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
			input:   sId,
			wantErr: nil,
			want: &types.Status{
				Id:        sId,
				Name:      "Existing project",
				ProjectId: project.Id,
			},
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", tt.want.Id, tt.want.Name, tt.want.ProjectId)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			ctx := context.Background()
			got, err := s.GetStatusById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetStatusById() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Status{}, "CreatedAt", "UpdatedAt", "Project")); diff != "" {
				t.Fatalf("GetStatusById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetStatusesByOwnerId(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	project := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "Discription of the project",
		OwnerId:     owner.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   string
		want    []types.Status
	}{
		"invalid owner id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
			want:    []types.Status{},
		},
		"2 statuses": {
			input:   project.Id,
			wantErr: nil,
			want: []types.Status{
				{
					Id:        uuid.NewString(),
					Name:      "Status 1",
					ProjectId: project.Id,
				},
				{
					Id:        uuid.NewString(),
					Name:      "Status 2",
					ProjectId: project.Id,
				},
			},
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		for _, st := range tt.want {
			_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", st.Id, st.Name, st.ProjectId)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			ctx := context.Background()
			got, err := s.GetStatusesByProjectId(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetStatusesByProjectId() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Status{}, "CreatedAt", "UpdatedAt", "Project")); diff != "" {
				t.Fatalf("GetStatusesByProjectId() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAddStatus(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	project := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "Discription of the project",
		OwnerId:     owner.Id,
	}
	tests := map[string]struct {
		input   *AddStatusInput
		wantErr error
	}{
		"invalid project id": {
			input: &AddStatusInput{
				Name:      "project",
				ProjectId: "invalid-id",
			},
			wantErr: ErrFailedValidation,
		},
		"invalid name": {
			input: &AddStatusInput{
				ProjectId: project.Id,
			},
			wantErr: ErrFailedValidation,
		},
		"non-existent project id": {
			input: &AddStatusInput{
				Name:      "status",
				ProjectId: uuid.NewString(),
			},
			wantErr: ErrInternal,
		},
		"sucsessfull add": {
			input: &AddStatusInput{
				Name:      "status",
				ProjectId: project.Id,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			ctx := context.Background()
			err := s.AddStatus(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("AddStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	project := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "Discription of the project",
		OwnerId:     owner.Id,
	}
	status := types.Status{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   *UpdateStatusInput
	}{
		"invalid id": {
			input: &UpdateStatusInput{
				Id:   "invalid-id",
				Name: "New status name",
			},
			wantErr: ErrFailedValidation,
		},
		"non-existent status id": {
			input: &UpdateStatusInput{
				Id:   uuid.NewString(),
				Name: "New status name",
			},
			wantErr: ErrFailedToUpdate,
		},
		"sucsessfull update": {
			input: &UpdateStatusInput{
				Id:   status.Id,
				Name: "New status name",
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			ctx := context.Background()
			err := s.UpdateStatus(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("UpdateStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteStatusById(t *testing.T) {
	s, cleanup := setupService(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	project := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "Discription of the project",
		OwnerId:     owner.Id,
	}
	status := types.Status{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   string
	}{
		"invalid id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
		},
		"non-existent status id": {
			input:   uuid.NewString(),
			wantErr: ErrFailedToUpdate,
		},
		"sucsessfull delete": {
			input:   status.Id,
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			ctx := context.Background()
			err := s.DeleteStatusById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("UpdateStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
