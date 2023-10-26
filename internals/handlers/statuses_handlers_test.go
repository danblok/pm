package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/danblok/pm/internals/service"
	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestHandleHandleGetStatusById(t *testing.T) {
	app, cleanup := setupApp(t)

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
		wantCode int
		want     *types.Status
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
			input:    sId,
			wantCode: http.StatusOK,
			want: &types.Status{
				Id:        sId,
				Name:      "Status 1",
				ProjectId: project.Id,
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
		if tt.want != nil {
			_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", tt.want.Id, tt.want.Name, tt.want.ProjectId)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleGetStatusById(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetStatusById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleGetStatusesByOwnerId(t *testing.T) {
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
	tests := map[string]struct {
		wantCode int
		input    string
		want     []types.Status
	}{
		"invalid owner id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
			want:     []types.Status{},
		},
		"2 projects": {
			input:    project.Id,
			wantCode: http.StatusOK,
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
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", project.Id, project.Name, project.Description, project.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		for _, st := range tt.want {
			_, err = app.Service.DB.Exec("INSERT INTO statuses (id, name, project_id) VALUES ($1, $2, $3)", st.Id, st.Name, st.ProjectId)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			q := make(url.Values)
			q.Set("pid", tt.input)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandleGetStatusesByOwner(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetStatusesByOwner() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleAddStatus(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Name      string `json:"name"`
		ProjectId string `json:"project_id"`
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
	tests := map[string]struct {
		input    *input
		wantCode int
	}{
		"invalid owner id": {
			input: &input{
				Name:      "project",
				ProjectId: "invalid-id",
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid name": {
			input: &input{
				ProjectId: owner.Id,
			},
			wantCode: http.StatusBadRequest,
		},
		"non-existent owner id": {
			input: &input{
				Name:      "project",
				ProjectId: uuid.NewString(),
			},
			wantCode: http.StatusInternalServerError,
		},
		"sucsessfull add": {
			input: &input{
				Name:      "project",
				ProjectId: project.Id,
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
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts", "statuses"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandlePostStatus(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandlePostStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleUpdateStatus(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
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
	tests := map[string]struct {
		param    string
		wantCode int
		input    *input
	}{
		"invalid id": {
			param: "invalid-id",
			input: &input{
				Name: "New project name",
			},
			wantCode: http.StatusBadRequest,
		},
		"non-existent project id": {
			param: uuid.NewString(),
			input: &input{
				Name: "New project name",
			},
			wantCode: http.StatusBadRequest,
		},
		"sucsessfull update": {
			param: status.Id,
			input: &input{
				Name: "New project name",
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
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts", "projects", "statuses"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.param)
			app.HandlePatchStatus(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleUpdateStatus(want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleDeleteStatusById(t *testing.T) {
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
	status := types.Status{
		Id:        uuid.NewString(),
		Name:      "in progress",
		ProjectId: project.Id,
	}
	tests := map[string]struct {
		wantCode int
		input    string
	}{
		"invalid id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
		},
		"non-existent project id": {
			input:    uuid.NewString(),
			wantCode: http.StatusBadRequest,
		},
		"sucsessfull delete": {
			input:    status.Id,
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

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts", "projects", "statuses"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleDeleteStatus(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleDeleteStatus() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
