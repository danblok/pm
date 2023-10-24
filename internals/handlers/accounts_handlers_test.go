package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danblok/pm/internals/service"
	"github.com/danblok/pm/internals/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestHandleGetAccount(t *testing.T) {
	app, cleanup := setupApp(t)

	accId := uuid.NewString()

	tests := map[string]struct {
		wantCode int
		input    string
		want     *types.Account
	}{
		"existent": {
			input:    accId,
			wantCode: http.StatusOK,
			want: &types.Account{
				Id:    accId,
				Name:  "username",
				Email: "username@test.com",
			},
		},
		"non-existent": {
			input:    uuid.NewString(),
			wantCode: http.StatusBadRequest,
			want:     nil,
		},
		"invalid id": {
			input:    "invald-id",
			wantCode: http.StatusBadRequest,
			want:     nil,
		},
	}

	for name, tt := range tests {
		if tt.want != nil {
			_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", tt.want.Id, tt.want.Email, tt.want.Name)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleGetAccount(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetAccount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleGetAllAccounts(t *testing.T) {
	app, cleanup := setupApp(t)

	tests := map[string]struct {
		wantCode int
		want     []types.Account
	}{
		"2 accounts": {
			want: []types.Account{
				{
					Id:    uuid.NewString(),
					Name:  "username 1",
					Email: "username1@test.com",
				},
				{
					Id:    uuid.NewString(),
					Name:  "username 2",
					Email: "username2@test.com",
				},
			},
			wantCode: http.StatusOK,
		},
		"0 accounts": {
			wantCode: http.StatusOK,
			want:     []types.Account{},
		},
	}

	for name, tt := range tests {
		for _, acc := range tt.want {
			_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)
			if err != nil {
				t.Fatal(service.ErrFailedToPrepareTest, err)
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandleGetAllAccounts(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleGetAllAccounts() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandlePostAccount(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar,omitempty"`
	}
	tests := map[string]struct {
		wantCode int
		input    *input
	}{
		"succsessfull add": {
			input: &input{
				Name:  "username",
				Email: "username@test.com",
			},
			wantCode: http.StatusCreated,
		},
		"invalid name": {
			input: &input{
				Name:  "",
				Email: "username@test.com",
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid email": {
			input: &input{
				Name:  "Project",
				Email: "",
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for name, tt := range tests {
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			app.HandlePostAccount(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleUpdateAccount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleUpdateAccount(t *testing.T) {
	app, cleanup := setupApp(t)

	type input struct {
		Name   string `json:"name,omitempty"`
		Email  string `json:"email,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}
	acc := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantCode int
		input    *input
		param    string
	}{
		"succsessfull update": {
			param: acc.Id,
			input: &input{
				Name:  "New project",
				Email: "newusername@test.com",
			},
			wantCode: http.StatusOK,
		},
		"non-existent id": {
			param: uuid.NewString(),
			input: &input{
				Name:  "New project",
				Email: "newusername@test.com",
			},
			wantCode: http.StatusBadRequest,
		},
		"invalid id": {
			param: "invalid-id",
			input: &input{
				Name:  "",
				Email: "username@test.com",
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest)
		}
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.param)
			app.HandleUpdateAccount(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleUpdateAccount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleDeleteAccount(t *testing.T) {
	app, cleanup := setupApp(t)

	acc := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@test.com",
	}
	tests := map[string]struct {
		wantCode int
		input    string
	}{
		"succsessfull delete": {
			input:    acc.Id,
			wantCode: http.StatusOK,
		},
		"non-existent": {
			input:    uuid.NewString(),
			wantCode: http.StatusBadRequest,
		},
		"invalid id": {
			input:    "invalid-id",
			wantCode: http.StatusBadRequest,
		},
	}

	for name, tt := range tests {
		_, err := app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)
		if err != nil {
			t.Fatal(service.ErrFailedToPrepareTest, err)
		}
		t.Run(name, func(t *testing.T) {
			t.Cleanup(cleanup("accounts"))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.input)
			app.HandleUpdateAccount(c)

			gotCode := res.Code
			if diff := cmp.Diff(tt.wantCode, gotCode); diff != "" {
				t.Fatalf("HandleDeleteAccount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
