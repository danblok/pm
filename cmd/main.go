package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/danblok/pm/internals/handlers"
	"github.com/danblok/pm/internals/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

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

	api := e.Group("/api")
	api.GET("/accounts/:id", app.HandleGetAccount)
	api.GET("/accounts", app.HandleGetAllAccounts)
	api.POST("/accounts", app.HandlePostAccount)
	api.DELETE("/accounts/:id", app.HandleDeleteAccount)
	api.PATCH("/accounts/:id", app.HandleUpdateAccount)
	api.GET("/projects/:id", app.HandleGetProjectById)
	api.GET("/projects", app.HandleGetProjectsByOwner)
	api.POST("/projects", app.HandlePostProject)
	api.PATCH("/projects/:id", app.HandleUpdateProject)
	api.DELETE("/projects/:id", app.HandleDeleteAccount)
	api.GET("/statuses/:id", app.HandleGetStatusById)
	api.GET("/statuses", app.HandleGetStatusesByOwner)
	api.POST("/statuses", app.HandlePostStatus)
	api.PATCH("/statuses/:id", app.HandleUpdateStatus)
	api.DELETE("/statuses/:id", app.HandleDeleteAccount)

	app.Logger.Info("Server started on http://localhost:3000")
	e.Logger.Fatal(e.Start(":3000"))
}
