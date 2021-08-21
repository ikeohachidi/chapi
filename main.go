package main

import (
	"github.com/ikeohachidi/chapi-be/models"
	"github.com/ikeohachidi/chapi-be/routers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := models.Connect()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := routers.App{
				c,
				db,
			}

			return next(cc)
		}
	})

	e.Logger.Fatal(e.Start(":1333"))
}
