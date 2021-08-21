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

// GetProjects retrieves a users projects
func GetUserProjects(c echo.Context) error {
	app := c.(App)

	if app.User.Id == 0 {
		log.Error("error: user id doesn't exist")
		return c.JSON(http.StatusBadRequest, Response{"Couldn't get user projects", false})
	}

	projects, err := app.Db.GetUserProjects(app.User.Id)
	if err != nil {
		log.Errorf("error: couldn't retrieve user projects: %v", err)
		return c.JSON(http.StatusBadRequest, Response{"Couldn't get user projects", false})
	}

	return c.JSON(http.StatusOK, projects)
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
	errResponseText := "Couldn't delete project"

	projectId, err := strconv.Atoi(app.QueryParam("id"))

	if err != nil {
		log.Fatalf("couldn't convert id param: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	if app.User.Id == 0 {
		log.Error("error: user id doesn't exist")
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	err = app.Db.DeleteProject(uint(projectId), app.User.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{"Project deleted successfully", true})
}
