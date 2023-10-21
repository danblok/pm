package service

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupService(t *testing.T) (*Service, func(...string) func()) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL_TEST"))
	if err != nil {
		t.Fatalf("connection to db: %s", err)
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
	return &Service{DB: db}, cleanup
}
