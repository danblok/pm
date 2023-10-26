package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	_ "github.com/danblok/pm/docs"
	"github.com/danblok/pm/internals/handlers"
	"github.com/danblok/pm/internals/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger PM API
// @version 1.0
// @description This is a PM server.

// @host localhost:3000
// @BasePath /api/v1
func main() {
	godotenv.Load(".env")
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		log.Fatal("POSTGRES_URL isn't specified")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Couldn't open connection to db: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Couldn't ping to db: ", err)
	}

	app := &handlers.App{
		Service: &service.Service{
			DB: db,
		},
		Logger: slog.Default(),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")
	api.GET("/accounts/:id", app.HandleGetAccount)
	api.GET("/accounts", app.HandleGetAllAccounts)
	api.POST("/accounts", app.HandlePostAccount)
	api.PATCH("/accounts/:id", app.HandlePatchAccount)
	api.DELETE("/accounts/:id", app.HandleDeleteAccount)
	api.GET("/projects/:id", app.HandleGetProjectById)
	api.GET("/projects", app.HandleGetProjectsByOwner)
	api.POST("/projects", app.HandlePostProject)
	api.PATCH("/projects/:id", app.HandlePatchProject)
	api.DELETE("/projects/:id", app.HandleDeleteAccount)
	api.GET("/statuses/:id", app.HandleGetStatusById)
	api.GET("/statuses", app.HandleGetStatusesByOwner)
	api.POST("/statuses", app.HandlePostStatus)
	api.PATCH("/statuses/:id", app.HandlePatchStatus)
	api.DELETE("/statuses/:id", app.HandleDeleteAccount)
	api.GET("/tasks/:id", app.HandleGetTaskById)
	api.GET("/tasks", app.HandleGetTasks)
	api.POST("/tasks", app.HandlePostTask)
	api.PATCH("/tasks/:id", app.HandlePatchTask)
	api.DELETE("/tasks/:id", app.HandleDeleteAccount)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	app.Logger.Info("Server started on http://localhost:3000")
	e.Logger.Fatal(e.Start(":3000"))
}
