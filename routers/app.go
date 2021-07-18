package routers

import (
	"github.com/ikeohachidi/chapi-be/models"
	"github.com/labstack/echo/v4"
)

type App struct {
	echo.Context
	Db   models.Conn
	Lama string
}

type Response struct {
	Data       interface{} `json:"message"`
	Successful bool        `json:"successful"`
}
