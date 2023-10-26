package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

// HandleGetStatus returns status
//
//	@Summary	Returns a status
//	@Tags		status
//	@Produce	json
//	@Param		id	path		string	true	"Status ID"
//	@Success	200	{object}	types.Status
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/statuses/{id}  [get]
func (a *App) HandleGetStatusById(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	s, err := a.Service.GetStatusById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.JSONPretty(http.StatusOK, &s, "  ")
}

// HandleGetStatusesByOwner lists all statuses of a project
//
//	@Summary	Returns all statuses of a project
//	@Tags		statuses
//	@Produce	json
//	@Param		pid	path	string	true	"Account ID"
//	@Success	200	{array}	types.Status
//	@Failure	400
//	@Failure	500
//	@Router		/statuses [get]
func (a *App) HandleGetStatusesByOwner(c echo.Context) error {
	ctx := c.Request().Context()
	pId := c.QueryParam("pid")

	sts, err := a.Service.GetStatusesByProjectId(ctx, pId)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.JSON(http.StatusOK, sts)
}

// HandlePostStatus creates a new status
//
//	@Summary	Create a new status
//	@Tags		status
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddStatusInput	true	"object of type AddStatusInput"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/statuses [post]
func (a *App) HandlePostStatus(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.AddStatusInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.AddStatus(ctx, input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusCreated)
}

// HandlePatchStatus patches an status
//
//	@Summary	Patche a status
//	@Tags		status
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddStatusInput	false	"body of type AddStatusInput"
//	@Param		id		path	string					true	"Status ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/statuses/{id} [patch]
func (a *App) HandlePatchStatus(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.UpdateStatusInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.UpdateStatus(ctx, input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}

// HandleDeleteStatus deletes an status
//
//	@Summary	Delete a status
//	@Tags		status
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Status ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/statuses/{id} [delete]
func (a *App) HandleDeleteStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := a.Service.DeleteStatusById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}
