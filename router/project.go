package router

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func Index(c echo.Context) error {
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

func CreateProjects(c echo.Context) error {
	app := c.(App)
	name := app.QueryParam("name")

	if name == "" {
		return c.JSON(http.StatusBadGateway, Response{"Bad Request", false})
	}

	id, err := app.Db.CreateProject(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Data:       "couldn't save project",
			Successful: false,
		})
	}

	return c.JSON(http.StatusOK, Response{id, true})
}

func ListProjects(c echo.Context) error {
	app := c.(App)

	projects, err := app.Db.ListProjects()
	if err != nil {
		return c.JSON(http.StatusBadGateway, Response{"Couldn't list projects", false})
	}

	return c.JSON(http.StatusOK, Response{projects, true})
}

func DeleteProject(c echo.Context) error {
	app := c.(App)

	projectId, err := strconv.Atoi(app.QueryParam("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{"Couldn't delete project", false})
	}

	err = app.Db.DeleteProject(uint(projectId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{"Couldn't delete project", false})
	}

	return c.JSON(http.StatusOK, Response{"Project deleted successfully", true})
}
