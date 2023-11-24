package endpoint

import (
	"go-Tsv/internal/config"
	"go-Tsv/internal/database/databaseStruct"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	cfg *config.Config

	svc DatabaseService
}

type DatabaseService interface {
	GetData(int, int) ([]databaseStruct.Messages, error)
}

func New(cfg *config.Config, db DatabaseService) *Endpoint {
	return &Endpoint{
		cfg: cfg,
		svc: db,
	}
}

func (e *Endpoint) HandlerGetData(ctx echo.Context) error {
	params := ctx.QueryParams()
	page, pageSize, err := e.parseParam(params)
	if err != nil {
		e.cfg.ErrorLog.Println("func HandlerGetData: ", err)
		return ctx.JSON(http.StatusOK, err.Error())
	}

	result, err := e.svc.GetData(page, pageSize)
	if err != nil {
		e.cfg.ErrorLog.Println("func HandlerGetData: ", err)
		return ctx.JSON(http.StatusOK, err.Error())
	}

	return ctx.JSON(http.StatusOK, result)
}

func (e *Endpoint) parseParam(p url.Values) (int, int, error) {
	pageSize := -1
	page := -1
	var err error

	if p.Get("page") != "" {
		page, err = strconv.Atoi(p.Get("page"))
		if err != nil {
			return 0, 0, err
		}
	}
	if p.Get("pageSize") != "" {
		pageSize, err = strconv.Atoi(p.Get("pageSize"))
		if err != nil {
			return 0, 0, err
		}
	}

	return page, pageSize, nil
}
