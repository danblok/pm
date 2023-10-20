package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Test data
var (
	addUserJson    = `{"name":"handlers_username", "email": "handlers_username@gmail.com"}`
	updateUserJson = `{"name":"updated_username"}`
)

func TestHandleGetAccount(t *testing.T) {
	app, cleanUp := setupApp(t)
	defer cleanUp("accounts")

	acc := types.Account{
		Id:    uuid.NewString(),
		Email: "username@gmail.com",
		Name:  "username",
	}
	app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/accounts/:id")
	c.SetParamNames("id")
	c.SetParamValues(acc.Id)

	err := app.HandleGetAccount(c)
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusOK
	got := res.Code
	if want != got {
		t.Fatalf("want: %d, recieved: %d", want, got)
	}
}

func TestHandleGetAllAccounts(t *testing.T) {
	app, cleanUp := setupApp(t)
	defer cleanUp("accounts")

	accs := []types.Account{
		{
			Id:    uuid.NewString(),
			Email: "username1@gmail.com",
			Name:  "username1",
		},
		{
			Id:    uuid.NewString(),
			Email: "username2@gmail.com",
			Name:  "username2",
		},
	}
	query := "INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3), ($4, $5, $6)"
	app.Service.DB.Exec(query, accs[0].Id, accs[0].Email, accs[0].Name)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	err := app.HandleGetAllAccounts(c)
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusOK
	got := res.Code
	if want != got {
		t.Fatalf("want: %d, recieved: %d", want, got)
	}
}

func TestHandlePostAccount(t *testing.T) {
	app, cleanUp := setupApp(t)
	defer cleanUp("accounts")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(addUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	err := app.HandlePostAccount(c)
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusCreated
	got := res.Code
	if want != got {
		t.Fatalf("want: %d, recieved: %d", want, got)
	}
}

func TestHandlePatchAccount(t *testing.T) {
	app, cleanUp := setupApp(t)
	defer cleanUp("accounts")
	defer func() {
		_, err := app.Service.DB.Exec("DELETE FROM accounts")
		if err != nil {
			t.Fatal("couldn't clean up the accounts table", err)
		}
	}()

	acc := types.Account{
		Id:    uuid.NewString(),
		Email: "username@gmail.com",
		Name:  "username",
	}
	app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(updateUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(acc.Id)

	err := app.HandleUpdateAccount(c)
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusOK
	got := res.Code
	if want != got {
		t.Fatalf("want: %d, recieved: %d", want, got)
	}
}

func TestHandleDeleteAccount(t *testing.T) {
	app, cleanUp := setupApp(t)
	defer cleanUp("accounts")

	acc := types.Account{
		Id:    uuid.NewString(),
		Email: "username@gmail.com",
		Name:  "username",
	}
	app.Service.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(acc.Id)

	err := app.HandleDeleteAccount(c)
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusOK
	got := res.Code
	if want != got {
		t.Fatalf("want: %d, recieved: %d", want, got)
	}
}
