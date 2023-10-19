package service

import (
	"database/sql"
	"errors"
)

var (
	ErrFailedValidation  = errors.New("failed validation")
	ErrFailedToUpdate    = errors.New("failed validation")
	ErrFaildeTestPrepare = errors.New("failed to prepare test data")
)

type Service struct {
	DB *sql.DB
}
