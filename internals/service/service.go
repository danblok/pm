package service

import (
	"database/sql"
	"errors"
)

var (
	ErrFailedValidation    = errors.New("failed validation")
	ErrFailedToUpdate      = errors.New("failed to update data")
	ErrFailedToPrepareTest = errors.New("failed to prepare test")
	ErrFailedToInsert      = errors.New("failed to insert data")
	ErrInternal            = errors.New("failed internal")
	ErrNotFound            = errors.New("not found")
)

type Service struct {
	DB *sql.DB
}
