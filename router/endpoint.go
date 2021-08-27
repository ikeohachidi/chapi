package router

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func HandleFrontend(c echo.Context) error {
	tmpl, err := template.New("index").ParseFiles("../frontend/index.html")
	if err != nil {
		log.Warnf("couldn't parse index.html file: %v", err)
	}

	err = tmpl.Execute(c.Response().Writer, nil)
	if err != nil {
		log.Warnf("couldn't exectute template: %v", err)
	}

	return nil
}

func StartProxy(c echo.Context) error {
	// blog.chapi.com/openmap
	app := c.(App)

	splitHost := strings.Split(c.Request().Host, ".")
	if splitHost[0] == "chapi" || len(splitHost) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	route, err := app.Db.GetRouteFromNameAndPath(splitHost[0], c.Request().URL.Path)
	if err != nil {
		log.Errorf("error getting project: %v", err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	return nil
}

func RunFrontendOrProxy(c echo.Context) error {
	host := c.Request().Host

	splitHost := strings.Split(host, ".")

	if splitHost[0] == "chapi" {
		HandleFrontend(c)
		return nil
	}

	StartProxy(c)

	return nil
}
