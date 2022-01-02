package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SaveProject(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't save project"

	project := model.Project{
		UserID: app.User.ID,
	}

	err := json.NewDecoder(c.Request().Body).Decode(&project)
	if err != nil {
		log.Errorf("couldn't decode request body: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	if c.Request().Method == http.MethodPost {
		err := project.Create(app.Conn.Db)
		if err != nil {
			log.Errorf("error creating project: %v", err)

			return c.JSON(http.StatusInternalServerError, Response{
				Data:       errResponseText,
				Successful: false,
			})
		}
	}
	if c.Request().Method == http.MethodPut {
		err := project.Update(app.Conn.Db)
		if err != nil {
			log.Errorf("error creating project: %v", err)

			return c.JSON(http.StatusInternalServerError, Response{
				Data:       errResponseText,
				Successful: false,
			})
		}
	}

	return c.JSON(http.StatusOK, Response{project, true})
}

func DoesProjectExist(c echo.Context) error {
	app := c.(App)
	project := model.Project{
		Name: c.QueryParam("name"),
	}
	errResponseText := "couldn't get result"

	projectExists, err := project.ProjectExists(app.Conn.Db)
	if err != nil {
		log.Errorf("error occured checking if project exists: %v", err)
		return c.JSON(http.StatusInternalServerError, Response{errResponseText, false})
	}

	if !projectExists {
		return c.JSON(http.StatusOK, Response{false, true})
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
	project := model.Project{
		UserID: app.User.ID,
	}

	projects, err := project.GetUserProjects(app.Conn.Db)
	if err != nil {
		log.Errorf("couldn't retrieve user projects: %v", err)
		return c.JSON(http.StatusBadRequest, Response{"Couldn't get user projects", false})
	}

	return c.JSON(http.StatusOK, Response{projects, true})
}

func ListProjects(c echo.Context) error {
	app := c.(App)

	projects, err := model.ListProjects(app.Conn.Db)
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
		log.Errorf("couldn't convert id param: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	if app.User.ID == 0 {
		log.Error("user id doesn't exist")
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	project := model.Project{
		ID:     uint(projectID),
		UserID: uint(app.User.ID),
	}

	err = project.DeleteProject(app.Conn.Db)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{"Project deleted successfully", true})
}
