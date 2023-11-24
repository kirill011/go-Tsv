package api

import (
	"fmt"
	"go-Tsv/internal/api/endpoint"
	"go-Tsv/internal/config"
	"go-Tsv/internal/database"

	"github.com/labstack/echo/v4"
)

type Api struct {
	db   *database.Database
	e    *endpoint.Endpoint
	cfg  *config.Config
	echo *echo.Echo
}

func New(db *database.Database, cfg *config.Config) *Api {
	a := &Api{
		db:   db,
		cfg:  cfg,
		e:    endpoint.New(cfg, db),
		echo: echo.New(),
	}

	a.echo.GET("/getData", a.e.HandlerGetData)

	return a
}

func (a *Api) Run() error {
	a.cfg.InfoLog.Println("Server running")

	err := a.echo.Start(fmt.Sprintf(":%d", a.cfg.ApiPort))
	if err != nil {
		a.echo.Logger.Fatal(err)
	}

	return nil
}
