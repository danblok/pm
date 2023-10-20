package service

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func setupServiceLifetime(t *testing.T) (*Service, func(...string)) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL_TEST"))
	if err != nil {
		t.Fatalf("connection to db: %s", err)
	}
	return &Service{DB: db}, func(tables ...string) {
		for _, table := range tables {
			_, err = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
			if err != nil {
				t.Fatal(err)
			}
		}
		db.Close()
	}
}
