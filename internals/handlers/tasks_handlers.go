package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/danblok/pm/internals/types"
	"github.com/labstack/echo/v4"
)

// HandleGetTask returns task
//
//	@Summary	Returns a task
//	@Tags		task
//	@Produce	json
//	@Param		id	path		string	true	"Task ID"
//	@Success	200	{object}	types.Task
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/tasks/{id} [get]
func (a *App) HandleGetTaskById(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	p, err := a.Service.GetTaskById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.JSON(http.StatusOK, &p)
}

// HandleGetTasks lists tasks of a project
//
//	@Summary	Returns list of tasks of a project
//	@Tags		tasks
//	@Produce	json
//	@Param		pid	path	string	true	"Project ID"
//	@Param		sid	path	string	false	"Status ID"
//	@Success	200	{array}	types.Task
//	@Failure	400
//	@Failure	500
//	@Router		/tasks [get]
func (a *App) HandleGetTasks(c echo.Context) error {
	ctx := c.Request().Context()
	pId := c.QueryParam("pid")
	sId := c.QueryParam("sid")
	var tks []types.Task
	var err error
	if sId != "" {
		tks, err = a.Service.GetTasksOfProjectByStatusId(ctx, pId, sId)
		if err != nil {
			return a.UnwrapError(c, "", err)
		}
	} else {
		tks, err = a.Service.GetTasksByProjectId(ctx, pId)
		if err != nil {
			return a.UnwrapError(c, "", err)
		}
	}

	return c.JSON(http.StatusOK, tks)
}

// HandlePostTask creates a new task
//
//	@Summary	Create a new task
//	@Tags		task
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddTaskInput	true	"object of type AddTaskInput"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/tasks [post]
func (a *App) HandlePostTask(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.AddTaskInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.AddTask(ctx, input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusCreated)
}

// HandlePatchTask patches a task
//
//	@Summary	Patche a task
//	@Tags		task
//	@Accept		json
//	@Produce	json
//	@Param		body	body	service.AddTaskInput	false	"body of type AddTaskInput"
//	@Param		id		path	string					true	"Task ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/tasks/{id} [patch]
func (a *App) HandlePatchTask(c echo.Context) error {
	ctx := c.Request().Context()
	input := new(service.UpdateTaskInput)
	err := c.Bind(input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	err = a.Service.UpdateTask(ctx, input)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}

// HandleDeleteTask deletes a task
//
//	@Summary	Delete a task
//	@Tags		task
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Task ID"
//	@Success	200
//	@Failure	400	{object}	types.HTTPError
//	@Failure	404	{object}	types.HTTPError
//	@Failure	500	{object}	types.HTTPError
//	@Router		/tasks/{id} [delete]
func (a *App) HandleDeleteTask(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := a.Service.DeleteTaskById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}
