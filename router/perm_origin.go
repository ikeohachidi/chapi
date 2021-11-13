package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi-be/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SavePermOrigins(c echo.Context) error {
	app := c.(App)

	var permOrigin model.PermOrigin
	var err error

	err = json.NewDecoder(app.Request().Body).Decode(&permOrigin)
	if err != nil {
		log.Errorf("couldn't create permission origin: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if c.Request().Method == "POST" {
		err = app.Db.CreatePermOrigin(&permOrigin)
	}

	if c.Request().Method == "PUT" {
		err = app.Db.UpdatePermOrigin(permOrigin)
	}

	if err != nil {
		log.Errorf("couldn't save permission origin: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, Response{permOrigin, true})
}

// GetPermOrigins will get perm_origin for a single route
func GetPermOrigins(c echo.Context) error {
	app := c.(App)

	routeID, err := strconv.Atoi(c.QueryParam("route_id"))
	if err != nil {
		log.Errorf("couldn't get permisson origins %v", err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	permOrigins, err := app.Db.GetPermOrigins(uint(routeID))
	if err != nil {
		log.Errorf("couldn't get permisson origins %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, Response{permOrigins, true})
}

func DeletePermOrigin(c echo.Context) error {
	app := c.(App)

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		log.Errorf("error converting query Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	routeID, err := strconv.Atoi(c.QueryParam("route_id"))
	if err != nil {
		log.Errorf("error converting route Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	permOrigin := model.PermOrigin{
		ID:      uint(id),
		RouteID: uint(routeID),
	}

	err = app.Db.DeletePermOrigin(permOrigin)
	if err != nil {
		log.Fatalf("couldn't delete permission origin %v", err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
