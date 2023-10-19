package handlers

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/danblok/pm/internals/service"
	_ "github.com/lib/pq"
)

func setupApp(t *testing.T) *App {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL_TEST"))
	if err != nil {
		t.Fatal("db connection err: ", err)
		return nil
	}
	return &App{
		Service: &service.Service{
			DB: db,
		},
		Logger: slog.Default(),
	}
}
