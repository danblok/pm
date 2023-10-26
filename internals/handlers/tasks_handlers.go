package handlers

import (
	"net/http"

	"github.com/danblok/pm/internals/service"
	"github.com/danblok/pm/internals/types"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetTaskById(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	p, err := a.Service.GetTaskById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.JSON(http.StatusOK, &p)
}

func (a *App) HandleGetTasks(c echo.Context) error {
	ctx := c.Request().Context()
	pId := c.QueryParam("pid")
	var tks []types.Task
	var err error
	if pId != "" {
		tks, err = a.Service.GetTasksByProjectId(ctx, pId)
		if err != nil {
			return a.UnwrapError(c, "", err)
		}
	} else {
		sId := c.QueryParam("sid")
		tks, err = a.Service.GetTasksByStatusId(ctx, sId)
		if err != nil {
			return a.UnwrapError(c, "", err)
		}
	}

	return c.JSON(http.StatusOK, tks)
}

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

func (a *App) HandleUpdateTask(c echo.Context) error {
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

func (a *App) HandleDeleteTask(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := a.Service.DeleteTaskById(ctx, id)
	if err != nil {
		return a.UnwrapError(c, "", err)
	}

	return c.NoContent(http.StatusOK)
}
