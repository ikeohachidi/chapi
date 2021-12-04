package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func SaveQuery(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't save query"

	query := model.Query{
		UserID: app.User.ID,
	}

	err := json.NewDecoder(c.Request().Body).Decode(&query)
	if err != nil {
		log.Errorf("couldn't decode query request body: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	HTTPMethod := c.Request().Method

	if HTTPMethod == http.MethodPost {
		err = query.Create(app.Conn.Db)
	}

	if HTTPMethod == http.MethodPut {
		err = query.Update(app.Conn.Db)
	}

	if err != nil {
		log.Errorf("error saving query to db: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return sendOkResponse(c, query)
}

func GetQueries(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't get route queries"

	routeID, err := strconv.Atoi(c.QueryParam("route"))
	if err != nil {
		log.Errorf("error converting route Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	if app.User.ID == 0 || routeID == 0 {
		log.Errorf("error getting user id or route id:\n userID: %v \n routeID: %v")
		return sendErrorResponse(c, http.StatusBadRequest, errResponseText)
	}

	query := model.Query{
		RouteID: uint(routeID),
		UserID:  app.User.ID,
	}

	queries, err := query.GetRouteQueries(app.Conn.Db)
	if err != nil {
		log.Errorf("error getting route queries from db: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return sendOkResponse(c, queries)
}

func DeleteQuery(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't delete route query"

	queryID, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		log.Errorf("error converting query Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	routeID, err := strconv.Atoi(c.QueryParam("route_id"))
	if err != nil {
		log.Errorf("error converting route Id string to int: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	var query = model.Query{
		UserID:  app.User.ID,
		ID:      uint(queryID),
		RouteID: uint(routeID),
	}

	err = query.Delete(app.Conn.Db)
	if err != nil {
		log.Errorf("error running delete route query query: %v", err)
		return c.JSON(http.StatusBadRequest, Response{errResponseText, false})
	}

	return sendOkResponse(c, "query deleted successfully")
}
