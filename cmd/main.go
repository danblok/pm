package main

import (
	"database/sql"
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
	url := os.Getenv("POSTGRES_URL_ALT")
	if url == "" {
		log.Fatal("POSTGRES_URL_ALT isn't specified")
	}

	db, err := sql.Open("postgres", url)
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


	app.logger.Info("Server started on http://localhost:3000")
	e.Logger.Fatal(e.Start(":3000"))
}
