package service

import (
	"database/sql"
	"errors"
)

var ErrFailedValidation = errors.New("failed validation")

type Service struct {
	DB *sql.DB
}
