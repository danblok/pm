package service

import (
	"database/sql"
	"os"
	"testing"
)

func setupService(t *testing.T) *Service {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL_TEST"))
	if err != nil {
		t.Fatalf("connection to db: %s", err)
	}
	return &Service{DB: db}
}
