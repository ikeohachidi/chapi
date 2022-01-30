package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func GetMergeOptions(c echo.Context) error {
	app := c.(App)

	routeId, err := strconv.Atoi(c.QueryParam("route_id"))
	if err != nil {
		log.Errorf("couldn't get merge options: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	mergeOptions := model.MergeOptions{RouteId: uint(routeId)}

	if err = mergeOptions.GetRouteMergeOptions(app.Conn.Db); err != nil {
		log.Errorf("couldn't get merge options: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, mergeOptions)
}

func SaveMergeOptions(c echo.Context) error {
	app := c.(App)

	mergeOption := model.MergeOptions{}

	if err := json.NewDecoder(c.Request().Body).Decode(&mergeOption); err != nil {
		log.Errorf("couldn't update merge options: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := mergeOption.SaveMergeOptions(app.Conn.Db); err != nil {
		log.Errorf("couldn't update merge options: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
