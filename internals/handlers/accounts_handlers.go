package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

// HandleGetAccount returns account
//
//	@Summary	Returns an account by ID
//	@Tags		account
//	@Produce	json
//	@Param		id	path		string	true	"Account ID"
//	@Success	200	{object}	types.Account
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/accounts/{id} [get]
func (a *App) HandleGetAccount(c echo.Context) error {
	id := c.Param("id")
	acc, err := a.Service.GetAccountById(c.Request().Context(), id)
	if err != nil {
		return a.UnwrapError(c, "Service.GetAccountById error: ", err)
	}

	return c.JSON(http.StatusOK, acc)
}

// HandleGetAccounts lists all existing accounts
//
//	@Summary	Returns all accounts
//	@Tags		accounts
//	@Produce	json
//	@Success	200	{array}	types.Account
//	@Failure	400
//	@Failure	500
//	@Router		/accounts [get]
func (a *App) HandleGetAllAccounts(c echo.Context) error {
	accs, err := a.Service.GetAllAccounts(c.Request().Context())
	if err != nil {
		return a.UnwrapError(c, "Service.GetAllAccounts error: ", err)
	}

	return c.JSON(http.StatusOK, accs)
}

// HandlePostAccount creates a new account
//
//	@Summary	Create an account
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddAccountInput	true	"object of type AddAccountInput"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/accounts [post]
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

// HandlePatchAccount patches an account
//
//	@Summary	Patch an account
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddAccountInput	false	"body of type AddAccountInput"
//	@Param		id		path	string					true	"Account ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/accounts/{id} [patch]
func (a *App) HandlePatchAccount(c echo.Context) error {
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

// HandleDeleteAccount deletes an account
//
//	@Summary	Delete an account
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Account ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/accounts/{id} [delete]
func (a *App) HandleDeleteAccount(c echo.Context) error {
	id := c.Param("id")
	err := a.Service.DeleteAccountById(c.Request().Context(), id)
	if err != nil {
		return a.UnwrapError(c, "Service.DeleteAccount: ", err)
	}
	return c.NoContent(http.StatusOK)
}
