package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi-be/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SaveQuery(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't save query"

	var query model.Query

	err := json.NewDecoder(c.Request().Body).Decode(&query)
	if err != nil {
		log.Errorf("couldn't decode query request body: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	query.UserID = app.User.ID

	HTTPMethod := c.Request().Method

	if HTTPMethod == http.MethodPost {
		err = app.Db.SaveQuery(&query)
	}

	if HTTPMethod == http.MethodPut {
		err = app.Db.UpdateQuery(query)
	}

	if err != nil {
		log.Errorf("error saving query to db: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{query.ID, true})
}

func GetRouteQueries(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't get route queries"

	routeID, err := strconv.Atoi(c.Param("routeID"))
	if err != nil {
		log.Errorf("error converting route Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	queries, err := app.Db.GetRouteQueries(uint(routeID), app.User.ID)
	if err != nil {
		log.Errorf("error getting route queries from db: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, queries)
}

func DeleteRouteQuery(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't delete route query"

	queryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("error converting query Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	routeID, err := strconv.Atoi(c.QueryParam("routeId"))
	if err != nil {
		log.Errorf("error converting route Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	err = app.Db.DeleteQuery(uint(queryID), uint(routeID), app.User.ID)
	if err != nil {
		log.Errorf("error running delete route query query: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return c.JSON(http.StatusOK, Response{"query deleted successfully", true})
}
