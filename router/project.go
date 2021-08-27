package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func CreateProject(c echo.Context) error {
	app := c.(App)
	name := app.QueryParam("name")
	errResponseText := "couldn't save project"

	if name == "" {
		log.Error("name query parameter is not defined")
		return c.JSON(http.StatusBadGateway, Response{"Bad Request", false})
	}

	if app.User.ID == 0 {
		log.Error("user id doesn't exist")
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	id, err := app.Db.CreateProject(name, app.User.ID)
	if err != nil {
		log.Errorf("error creating project: %v", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Data:       errResponseText,
			Successful: false,
		})
	}

	return c.JSON(http.StatusOK, Response{id, true})
}

func CheckProjectExistence(c echo.Context) error {
	app := c.(App)
	userInput := c.QueryParam("userInput")
	errResponseText := "couldn't get result"

	projectExists, err := app.Db.ProjectExists(userInput)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{errResponseText, false})
	}

	if !projectExists {
		return c.JSON(http.StatusOK, Response{false, false})
	}

	return c.JSON(http.StatusOK, Response{true, true})
}

// GetProjects retrieves a users projects
func GetUserProjects(c echo.Context) error {
	app := c.(App)

	if app.User.ID == 0 {
		log.Error("user id doesn't exist")
		return c.JSON(http.StatusBadRequest, Response{"Couldn't get user projects", false})
	}

	projects, err := app.Db.GetUserProjects(app.User.ID)
	if err != nil {
		log.Errorf("couldn't retrieve user projects: %v", err)
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

	projectID, err := strconv.Atoi(app.Param("id"))

	if err != nil {
		log.Fatalf("couldn't convert id param: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	if app.User.ID == 0 {
		log.Error("user id doesn't exist")
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	err = app.Db.DeleteProject(uint(projectID), app.User.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{"Project deleted successfully", true})
}
