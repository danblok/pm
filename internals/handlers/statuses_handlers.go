package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetStatusById(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	s, err := a.Service.GetStatusById(ctx, id)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.JSONPretty(http.StatusOK, &s, "  ")
}

func (a *App) HandleGetStatusesByOwner(c echo.Context) error {
	ctx := c.Request().Context()
	pId := c.QueryParam("pid")

	sts, err := a.Service.GetStatusesByProjectId(ctx, pId)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.JSON(http.StatusOK, sts)
}

func (a *App) HandlePostStatus(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.AddStatusInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.AddStatus(ctx, input)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusCreated)
}

func (a *App) HandleUpdateStatus(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.UpdateStatusInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.UpdateStatus(ctx, input)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}

func (a *App) HandleDeleteStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := a.Service.DeleteStatusById(ctx, id)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}
