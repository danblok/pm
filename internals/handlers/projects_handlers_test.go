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

func TestHandleGetProjectById(t *testing.T) {
	app, cleanup := setupApp(t)

	pId := uuid.NewString()
	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantCode int
		want     *types.Project
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
			input:    pId,
			wantCode: http.StatusOK,
			want: &types.Project{
				Id:      pId,
				Name:    "Existing project",
				OwnerId: owner.Id,
			},
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", tt.want.Id, tt.want.Name, tt.want.Description, tt.want.OwnerId)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleGetProjectById(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetProjectById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleGetProjectsByOwnerId(t *testing.T) {
	app, cleanup := setupApp(t)

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantCode int
		input    string
		want     []types.Project
	}{
		"invalid owner id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
			want:     []types.Project{},
		},
		"2 projects": {
			input:    owner.Id,
			wantCode: http.StatusOK,
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
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		if tt.want != nil {
			for _, p := range tt.want {
				_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
				if err != nil {
					t.Fatal(service.ErrFailedToPrepareTest, err)
				}
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			q := make(url.Values)
			q.Set("oid", tt.input)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandleGetProjectsByOwner(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetProjectsByOwner() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleAddProject(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		OwnerId     string `json:"owner_id"`
	}
	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		input    *input
		wantCode int
	}{
		"invalid owner id": {
			input: &input{
				Name:    "project",
				OwnerId: "invalid-id",
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid name": {
			input: &input{
				OwnerId: owner.Id,
			},
			wantCode: http.StatusBadRequest,
		},
		"non-existent owner id": {
			input: &input{
				Name:    "project",
				OwnerId: uuid.NewString(),
			},
			wantCode: http.StatusInternalServerError,
		},
		"sucsessfull add": {
			input: &input{
				Name:    "project",
				OwnerId: owner.Id,
			},
			wantCode: http.StatusCreated,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandlePostProject(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandlePostProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleUpdateProject(t *testing.T) {
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
	p := types.Project{
		Id:      uuid.NewString(),
		Name:    "Project",
		OwnerId: owner.Id,
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
			param: p.Id,
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
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("projects", "accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.param)
			app.HandleUpdateProject(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleUpdateProject(want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleDeleteProjectById(t *testing.T) {
	app, cleanup := setupApp(t)

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
			input:    p.Id,
			wantCode: http.StatusOK,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		_, err = app.Service.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts", "projects"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleDeleteProject(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleDeleteProject() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
