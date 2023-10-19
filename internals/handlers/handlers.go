package handlers

import (
	"log/slog"

	"github.com/danblok/pm/internals/service"
)

type App struct {
	Service *service.Service
	Logger  *slog.Logger
}
