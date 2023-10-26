package service

import (
	"context"
	"testing"
	"time"

	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestGetTaskById(t *testing.T) {
	s, cleanup := setupService(t)

	tId := uuid.NewString()
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
	status := types.Task{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		wantErr error
		want    *types.Task
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
			input:   tId,
			wantErr: nil,
			want: &types.Task{
				Id:        tId,
				Name:      "Existing project",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().UTC(),
				End:       time.Now().UTC().AddDate(0, 0, 1),
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
		_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			_, err = s.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", tt.want.Id, tt.want.Name, tt.want.ProjectId, tt.want.StatusId, tt.want.Start, tt.want.End)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			ctx := context.Background()
			got, err := s.GetTaskById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetTaskById() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Task{}, "CreatedAt", "UpdatedAt", "Project", "Status"), cmpopts.EquateApproxTime(time.Millisecond)); diff != "" {
				t.Fatalf("GetTaskById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetTasksByProjectId(t *testing.T) {
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
	status := types.Task{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   string
		want    []types.Task
	}{
		"invalid owner id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
			want:    []types.Task{},
		},
		"2 statuses": {
			input:   project.Id,
			wantErr: nil,
			want: []types.Task{
				{
					Id:        uuid.NewString(),
					Name:      "Task 1",
					ProjectId: project.Id,
					StatusId:  status.Id,
					Start:     time.Now().UTC(),
					End:       time.Now().UTC(),
				},
				{
					Id:        uuid.NewString(),
					Name:      "Task 2",
					ProjectId: project.Id,
					StatusId:  status.Id,
					Start:     time.Now().UTC(),
					End:       time.Now().UTC(),
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
		_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		for _, ts := range tt.want {
			_, err = s.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", ts.Id, ts.Name, ts.ProjectId, ts.StatusId, ts.Start, ts.End)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			ctx := context.Background()
			got, err := s.GetTasksByProjectId(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetTasksByProjectId() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Task{}, "CreatedAt", "UpdatedAt", "Project"), cmpopts.EquateApproxTime(time.Millisecond)); diff != "" {
				t.Fatalf("GetTasksByProjectId() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetTasksOfProjectByStatusId(t *testing.T) {
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
	status := types.Task{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   struct{ statusId, projectId string }
		want    []types.Task
	}{
		"invalid status id": {
			input: struct {
				statusId  string
				projectId string
			}{statusId: "invalid-id", projectId: project.Id},
			wantErr: ErrFailedValidation,
			want:    []types.Task{},
		},
		"invalid project id": {
			input: struct {
				statusId  string
				projectId string
			}{projectId: "invalid-id", statusId: status.Id},
			wantErr: ErrFailedValidation,
			want:    []types.Task{},
		},
		"2 statuses": {
			input: struct {
				statusId  string
				projectId string
			}{statusId: status.Id, projectId: project.Id},
			wantErr: nil,
			want: []types.Task{
				{
					Id:        uuid.NewString(),
					Name:      "Task 1",
					ProjectId: project.Id,
					StatusId:  status.Id,
					Start:     time.Now().UTC(),
					End:       time.Now().UTC(),
				},
				{
					Id:        uuid.NewString(),
					Name:      "Task 2",
					ProjectId: project.Id,
					StatusId:  status.Id,
					Start:     time.Now().UTC(),
					End:       time.Now().UTC(),
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
		_, err = s.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}
		for _, ts := range tt.want {
			_, err = s.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", ts.Id, ts.Name, ts.ProjectId, ts.StatusId, ts.Start, ts.End)
			if err != nil {
				t.Fatal(ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			ctx := context.Background()
			got, err := s.GetTasksOfProjectByStatusId(ctx, tt.input.projectId, tt.input.statusId)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("GetTasksByStatusId() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(types.Task{}, "CreatedAt", "UpdatedAt", "Project"), cmpopts.EquateApproxTime(time.Millisecond)); diff != "" {
				t.Fatalf("GetTasksByStatusId() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
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
	status := types.Task{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		input   *AddTaskInput
		wantErr error
	}{
		"invalid project id": {
			input: &AddTaskInput{
				Name:      "task",
				ProjectId: "invalid-id",
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid status id": {
			input: &AddTaskInput{
				Name:      "task",
				StatusId:  "invalid-id",
				ProjectId: project.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid name": {
			input: &AddTaskInput{
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid start": {
			input: &AddTaskInput{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     "invalid start",
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid end": {
			input: &AddTaskInput{
				ProjectId: project.Id,
				Name:      "task",
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       "invalid end",
			},
			wantErr: ErrFailedValidation,
		},
		"end is less than start": {
			input: &AddTaskInput{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().AddDate(0, 0, 1).Format(time.DateTime),
				End:       time.Now().Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"non-existent project id": {
			input: &AddTaskInput{
				Name:      "task",
				ProjectId: uuid.NewString(),
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrInternal,
		},
		"non-existent status id": {
			input: &AddTaskInput{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  uuid.NewString(),
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrInternal,
		},
		"sucsessfull add": {
			input: &AddTaskInput{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
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
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			ctx := context.Background()
			err := s.AddTask(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("AddTask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
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
	task := types.Task{
		Id:        uuid.NewString(),
		Name:      "task",
		ProjectId: project.Id,
		StatusId:  status.Id,
	}
	tests := map[string]struct {
		wantErr error
		input   *UpdateTaskInput
	}{
		"invalid id": {
			input: &UpdateTaskInput{
				Id:       "invalid-id",
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid status id": {
			input: &UpdateTaskInput{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: "invalid-id",
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid start": {
			input: &UpdateTaskInput{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    "invalid start",
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"invalid end": {
			input: &UpdateTaskInput{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      "invalid end",
			},
			wantErr: ErrFailedValidation,
		},
		"end is less than start": {
			input: &UpdateTaskInput{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().AddDate(0, 0, 1).Format(time.DateTime),
				End:      time.Now().Format(time.DateTime),
			},
			wantErr: ErrFailedValidation,
		},
		"non-existent task id": {
			input: &UpdateTaskInput{
				Id:       uuid.NewString(),
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrFailedToUpdate,
		},
		"non-existent status id": {
			input: &UpdateTaskInput{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: uuid.NewString(),
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantErr: ErrInternal,
		},
		"sucsessfull update": {
			input: &UpdateTaskInput{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
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
		_, err = s.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", task.Id, task.Name, task.ProjectId, task.StatusId, task.Start, task.End)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			ctx := context.Background()
			err := s.UpdateTask(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("UpdateTask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteTaskById(t *testing.T) {
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
	status := types.Task{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	task := types.Task{
		Id:        uuid.NewString(),
		Name:      "task",
		ProjectId: project.Id,
		StatusId:  status.Id,
		Start:     time.Now().UTC(),
		End:       time.Now().AddDate(0, 0, 1).UTC(),
	}
	tests := map[string]struct {
		wantErr error
		input   string
	}{
		"invalid id": {
			input:   "invalid-id",
			wantErr: ErrFailedValidation,
		},
		"non-existent task id": {
			input:   uuid.NewString(),
			wantErr: ErrFailedToUpdate,
		},
		"sucsessfull delete": {
			input:   task.Id,
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
		_, err = s.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", task.Id, task.Name, task.ProjectId, task.StatusId, task.Start, task.End)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			ctx := context.Background()
			err := s.DeleteTaskById(ctx, tt.input)
			if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Fatalf("DeleteTaskById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
