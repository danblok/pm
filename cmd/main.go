package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/danblok/pm/internals/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type app struct {
	service *service.Service
	logger  *slog.Logger
}

func main() {
	godotenv.Load()
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		log.Fatal("POSTGRES_USER isn't specified")
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		log.Fatal("POSTGRES_PASSWORD isn't specified")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", user, password))
	if err != nil {
		log.Fatal("Couldn't open connection to db: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Couldn't ping to db: ", err)
	}

	app := &app{
		service: &service.Service{
			DB: db,
		},
		logger: slog.Default(),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: handlers setup

	app.logger.Info("Server started on http://localhost:3000")
	e.Logger.Fatal(e.Start(":3000"))
}
