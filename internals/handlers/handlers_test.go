package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/danblok/pm/internals/service"
	_ "github.com/lib/pq"
)

func setupApp(t *testing.T) (*App, func(...string) func()) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		t.Fatal("db connection err: ", err)
	}
	cleanup := func(tables ...string) func() {
		return func() {
			for _, table := range tables {
				_, err = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	return &App{
		Service: &service.Service{
			DB: db,
		},
		Logger: slog.Default(),
	}, cleanup
}
