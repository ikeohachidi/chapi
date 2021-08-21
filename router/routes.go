package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi-be/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SaveRoute(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't save route"
	var route model.Route

	err := json.NewDecoder(c.Request().Body).Decode(&route)

	if err != nil {
		log.Errorf("couldn't decode save route request body: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	routeId, err := app.Db.SaveRoute(route)
	if err != nil {
		log.Errorf("couldn't save route %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{routeId, true})
}

func GetProjectRoutes(c echo.Context) error {
	app := c.(App)
	projectId, _ := strconv.Atoi(c.Param("projectID"))
	errResponseText := "couldn't retrieve project route"

	routes, err := app.Db.GetRoutesByProjectId(uint(projectId))
	if err != nil {
		log.Errorf("couldn't retrieve project routes %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{routes, true})
}

func DeleteRoute(c echo.Context) error {
	app := c.(App)
	routeId, _ := strconv.Atoi(c.Param("id"))
	errResponseText := "couldn't delete route"

	if routeId == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err := app.Db.DeleteRoute(uint(routeId), app.User.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{errResponseText, false})
	}

	return nil
}
