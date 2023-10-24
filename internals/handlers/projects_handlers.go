package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetProjectById(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	p, err := a.Service.GetProjectById(ctx, id)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.JSONPretty(http.StatusOK, &p, "  ")
}

func (a *App) HandleGetProjectsByOwner(c echo.Context) error {
	ctx := c.Request().Context()
	pId := c.QueryParam("pid")

	pjs, err := a.Service.GetProjectsByOwnerId(ctx, pId)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.JSONPretty(http.StatusOK, pjs, "  ")
}

func (a *App) HandlePostProject(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.AddProjectInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.AddProject(ctx, input)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusCreated)
}

func (a *App) HandleUpdateProject(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.UpdateProjectInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.UpdateProject(ctx, input)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}

func (a *App) HandleDeleteProject(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := a.Service.DeleteProjectById(ctx, id)
	if err != nil {
		a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}
