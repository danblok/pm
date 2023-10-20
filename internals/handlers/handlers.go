package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

type App struct {
	Service *service.Service
	Logger  *slog.Logger
}

func (a *App) UnwrapError(c echo.Context, logMsg string, err error) error {
	a.Logger.Error(logMsg, err)
	if errors.Is(err, service.ErrInternal) {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusBadRequest)
}
