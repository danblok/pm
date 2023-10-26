package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/danblok/pm/internals/service"
	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestHandleGetTaskById(t *testing.T) {
	app, cleanup := setupApp(t)

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
		wantCode int
		want     *types.Task
		input    string
	}{
		"non-existent": {
			input:    uuid.NewString(),
			wantCode: http.StatusBadRequest,
			want:     nil,
		},
		"invalid id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
			want:     nil,
		},
		"existing": {
			input:    tId,
			wantCode: http.StatusOK,
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
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			_, err = app.Service.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", tt.want.Id, tt.want.Name, tt.want.ProjectId, tt.want.StatusId, tt.want.Start, tt.want.End)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleGetTaskById(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetTaskById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleGetTasksByProjectId(t *testing.T) {
	app, cleanup := setupApp(t)

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
		wantCode int
		input    string
		want     []types.Task
	}{
		"invalid task id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
			want:     []types.Task{},
		},
		"2 tasks": {
			input:    project.Id,
			wantCode: http.StatusOK,
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
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		for _, ts := range tt.want {
			_, err = app.Service.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", ts.Id, ts.Name, ts.ProjectId, ts.StatusId, ts.Start, ts.End)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			q := make(url.Values)
			q.Set("pid", tt.input)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandleGetTasks(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetTasksByProjectId() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleGetTasksByStatusId(t *testing.T) {
	app, cleanup := setupApp(t)

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
		wantCode int
		input    string
		want     []types.Task
	}{
		"invalid owner id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
			want:     []types.Task{},
		},
		"2 tasks": {
			input:    status.Id,
			wantCode: http.StatusOK,
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
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		for _, ts := range tt.want {
			_, err = app.Service.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", ts.Id, ts.Name, ts.ProjectId, ts.StatusId, ts.Start, ts.End)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			q := make(url.Values)
			q.Set("sid", tt.input)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandleGetTasks(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetTasksByStatusId() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleAddTask(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Start     string `json:"start"`
		End       string `json:"end"`
		Name      string `json:"name"`
		ProjectId string `json:"project_id"`
		StatusId  string `json:"status_id"`
	}
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
		input    *input
		wantCode int
	}{
		"invalid project id": {
			input: &input{
				Name:      "task",
				ProjectId: "invalid-id",
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid status id": {
			input: &input{
				Name:      "task",
				StatusId:  "invalid-id",
				ProjectId: project.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid name": {
			input: &input{
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid start": {
			input: &input{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     "invalid start",
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid end": {
			input: &input{
				ProjectId: project.Id,
				Name:      "task",
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       "invalid end",
			},
			wantCode: http.StatusBadRequest,
		},
		"end is less than start": {
			input: &input{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().AddDate(0, 0, 1).Format(time.DateTime),
				End:       time.Now().Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"non-existent project id": {
			input: &input{
				Name:      "task",
				ProjectId: uuid.NewString(),
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusInternalServerError,
		},
		"non-existent status id": {
			input: &input{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  uuid.NewString(),
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusInternalServerError,
		},
		"sucsessfull add": {
			input: &input{
				Name:      "task",
				ProjectId: project.Id,
				StatusId:  status.Id,
				Start:     time.Now().Format(time.DateTime),
				End:       time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusCreated,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandlePostTask(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandlePostTask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleUpdateTask(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Start    string `json:"start,omitempty"`
		End      string `json:"end,omitempty"`
		Id       string `param:"id"`
		Name     string `json:"name,omitempty"`
		StatusId string `json:"status_id"`
	}
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
		param    string
		wantCode int
		input    *input
	}{
		"invalid id": {
			input: &input{
				Id:       "invalid-id",
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid status id": {
			input: &input{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: "invalid-id",
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid start": {
			input: &input{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    "invalid start",
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid end": {
			input: &input{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      "invalid end",
			},
			wantCode: http.StatusBadRequest,
		},
		"end is less than start": {
			input: &input{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().AddDate(0, 0, 1).Format(time.DateTime),
				End:      time.Now().Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"non-existent task id": {
			input: &input{
				Id:       uuid.NewString(),
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusBadRequest,
		},
		"non-existent status id": {
			input: &input{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: uuid.NewString(),
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusInternalServerError,
		},
		"sucsessfull update": {
			input: &input{
				Id:       task.Id,
				Name:     "New task name",
				StatusId: status.Id,
				Start:    time.Now().Format(time.DateTime),
				End:      time.Now().AddDate(0, 0, 1).Format(time.DateTime),
			},
			wantCode: http.StatusOK,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", task.Id, task.Name, task.ProjectId, task.StatusId, task.Start, task.End)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.param)
			app.HandleUpdateTask(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleUpdateTask(want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleDeleteTaskById(t *testing.T) {
	app, cleanup := setupApp(t)

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
		wantCode int
		input    string
	}{
		"invalid id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
		},
		"non-existent task id": {
			input:    uuid.NewString(),
			wantCode: http.StatusBadRequest,
		},
		"sucsessfull delete": {
			input:    task.Id,
			wantCode: http.StatusOK,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", status.Id, status.Name, status.ProjectId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO tasks (id, name, project_id, status_id, \"start\", \"end\") VALUES ($1, $2, $3, $4, $5, $6)", task.Id, task.Name, task.ProjectId, task.StatusId, task.Start, task.End)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses", "tasks"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleDeleteTask(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleDeleteTask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
