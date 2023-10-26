package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/labstack/echo/v4"
)

// HandleGetProject returns project
//
//	@Summary	Returns a project
//	@Tags		project
//	@Produce	json
//	@Param		id	path		string	true	"Project ID"
//	@Success	200	{object}	types.Project
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/projects/{id} [get]
func (a *App) HandleGetProjectById(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	p, err := a.Service.GetProjectById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.JSONPretty(http.StatusOK, &p, "  ")
}

// HandleGetProjects lists all existing projects
//
//	@Summary	Returns all projects of an account
//	@Tags		projects
//	@Produce	json
//	@Param		pid	path	string	true	"Account ID"
//	@Success	200	{array}	types.Project
//	@Failure	400
//	@Failure	500
//	@Router		/projects [get]
func (a *App) HandleGetProjectsByOwner(c echo.Context) error {
	ctx := c.Request().Context()
	oId := c.QueryParam("oid")

	pjs, err := a.Service.GetProjectsByOwnerId(ctx, oId)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.JSONPretty(http.StatusOK, pjs, "  ")
}

// HandlePostProject creates a new project
//
//	@Summary	Create a new project
//	@Tags		project
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddProjectInput	true	"object of type AddProjectInput"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/projects [post]
func (a *App) HandlePostProject(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.AddProjectInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.AddProject(ctx, input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusCreated)
}

// HandlePatchProject patches an project
//
//	@Summary	Patche a project
//	@Tags		project
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddProjectInput	false	"body of type AddProjectInput"
//	@Param		id		path	string					true	"Project ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/projects/{id} [patch]
func (a *App) HandlePatchProject(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.UpdateProjectInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.UpdateProject(ctx, input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}

// HandleDeleteProject deletes an project
//
//	@Summary	Delete a project
//	@Tags		project
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Project ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/projects/{id} [delete]
func (a *App) HandleDeleteProject(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := a.Service.DeleteProjectById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}
