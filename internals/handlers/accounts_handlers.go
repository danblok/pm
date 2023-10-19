package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetAccount(c echo.Context) error {
	id := c.Param("id")
	acc, err := a.Service.GetAccountById(c.Request().Context(), id)
	if err != nil {
		a.Logger.Error("Service.GetAccountById error: ", err)
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			return echo.ErrBadRequest
		case errors.Is(err, sql.ErrNoRows):
			return echo.ErrBadRequest
		default:
			return echo.ErrInternalServerError
		}
	}

	return c.JSON(http.StatusOK, &acc)
}

func (a *App) HandleGetAllAccounts(c echo.Context) error {
	accs, err := a.Service.GetAllAccounts(c.Request().Context())
	if err != nil {
		a.Logger.Error("Service.GetAllAccounts error: ", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, accs)
}

func (a *App) HandlePostAccount(c echo.Context) error {
	var input service.AddAccountInput
	err := c.Bind(&input)
	if err != nil {
		a.Logger.Error("binding in HandlePostAccount input error: ", err)
		return echo.ErrBadRequest
	}

	err = a.Service.AddAccount(c.Request().Context(), &input)
	if err != nil {
		a.Logger.Error("Service.UpdateAccount error: ", err)
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			return echo.ErrBadRequest
		case errors.Is(err, service.ErrFailedToUpdate):
			return echo.ErrBadRequest
		default:
			return echo.ErrInternalServerError
		}
	}
	return c.NoContent(http.StatusCreated)
}

func (a *App) HandlePatchAccount(c echo.Context) error {
	var input service.UpdateAccountInput
	input.Id = c.Param("id")
	err := c.Bind(&input)
	if err != nil {
		a.Logger.Error("binding in HandlePatchAccount input error: ", err)
		return echo.ErrBadRequest
	}

	err = a.Service.UpdateAccount(c.Request().Context(), &input)
	if err != nil {
		a.Logger.Error("Service.UpdateAccount error: ", err)
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			return echo.ErrBadRequest
		case errors.Is(err, service.ErrFailedToUpdate):
			return echo.ErrBadRequest
		default:
			return echo.ErrInternalServerError
		}
	}
	return c.NoContent(http.StatusOK)
}

func (a *App) HandleDeleteAccount(c echo.Context) error {
	id := c.Param("id")
	err := a.Service.DeleteAccount(c.Request().Context(), id)
	if err != nil {
		a.Logger.Error("Service.DeleteAccount: ", err)
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			return echo.ErrBadRequest
		case errors.Is(err, service.ErrFailedToUpdate):
			return echo.ErrBadRequest
		default:
			return echo.ErrInternalServerError
		}
	}
	return c.NoContent(http.StatusOK)
}
