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

	var route = model.Route{
		UserID: app.User.ID,
	}

	err := json.NewDecoder(c.Request().Body).Decode(&route)

	if err != nil {
		log.Errorf("couldn't decode save route request body: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	HTTPMethod := c.Request().Method

	if HTTPMethod == http.MethodPost {
		err = app.Db.SaveRoute(&route)
	}

	if HTTPMethod == http.MethodPut {
		err = app.Db.UpdateRoute(route)
	}

	if err != nil {
		log.Errorf("couldn't save route %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{route.ID, true})
}

func GetProjectRoutes(c echo.Context) error {
	app := c.(App)
	projectID, _ := strconv.Atoi(c.Param("projectID"))
	errResponseText := "couldn't retrieve project route"

	if app.User.ID == 0 {
		log.Errorf("can't get project routes with invalid user id")
		return c.JSON(http.StatusBadRequest, nil)
	}

	routes, err := app.Db.GetRoutesByProjectId(uint(projectID), app.User.ID)
	if err != nil {
		log.Errorf("couldn't retrieve project routes %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{routes, true})
}

func DeleteRoute(c echo.Context) error {
	app := c.(App)
	routeID, _ := strconv.Atoi(c.Param("id"))
	errResponseText := "couldn't delete route"

	if routeID == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err := app.Db.DeleteRoute(uint(routeID), app.User.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{errResponseText, false})
	}

	return nil
}
