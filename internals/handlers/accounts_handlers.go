package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetAccount(c echo.Context) error {
	id := c.Param("id")
	acc, err := a.Service.GetAccountById(c.Request().Context(), id)
	if err != nil {
		return a.UnwrapError(c, "Service.GetAccountById error: ", err)
	}

	return c.JSON(http.StatusOK, acc)
}

func (a *App) HandleGetAllAccounts(c echo.Context) error {
	accs, err := a.Service.GetAllAccounts(c.Request().Context())
	if err != nil {
		return a.UnwrapError(c, "Service.GetAllAccounts error: ", err)
	}

	return c.JSON(http.StatusOK, accs)
}

func (a *App) HandlePostAccount(c echo.Context) error {
	var input service.AddAccountInput
	err := c.Bind(&input)
	if err != nil {
		return a.UnwrapError(c, "binding in HandlePostAccount input error: ", err)
	}

	err = a.Service.AddAccount(c.Request().Context(), &input)
	if err != nil {
		return a.UnwrapError(c, "Service.UpdateAccount error: ", err)
	}
	return c.NoContent(http.StatusCreated)
}

func (a *App) HandleUpdateAccount(c echo.Context) error {
	var input service.UpdateAccountInput
	input.Id = c.Param("id")
	err := c.Bind(&input)
	if err != nil {
		return a.UnwrapError(c, "binding in HandlePatchAccount input error: ", err)
	}

	err = a.Service.UpdateAccount(c.Request().Context(), &input)
	if err != nil {
		return a.UnwrapError(c, "Service.UpdateAccount error: ", err)
	}
	return c.NoContent(http.StatusOK)
}

func (a *App) HandleDeleteAccount(c echo.Context) error {
	id := c.Param("id")
	err := a.Service.DeleteAccountById(c.Request().Context(), id)
	if err != nil {
		return a.UnwrapError(c, "Service.DeleteAccount: ", err)
	}
	return c.NoContent(http.StatusOK)
}
