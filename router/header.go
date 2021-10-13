package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ikeohachidi/chapi-be/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func GetHeaders(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't retrieve header"

	routeID, err := strconv.Atoi(c.QueryParam("route"))
	if err != nil {
		log.Errorf("error converting route id query to int: %v", err)
		return sendErrorResponse(c, http.StatusInternalServerError, errResponseText)
	}

	if app.User.ID == 0 || routeID == 0 {
		log.Errorf("error getting user id or route id:\n userID: %v \n routeID: %v")
		return sendErrorResponse(c, http.StatusBadRequest, errResponseText)
	}

	header, err := app.Db.GetHeader(app.User.ID, uint(routeID))

	if err != nil {
		log.Errorf("error retrieving header: %v", err)
		return sendErrorResponse(c, http.StatusInternalServerError, errResponseText)
	}

	return sendOkResponse(c, header)
}

func SaveHeader(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't save header"

	var header model.Header

	routeID, err := strconv.Atoi(c.QueryParam("route"))
	if err != nil {
		log.Errorf("error converting route id query to int: %v", err)
		return sendErrorResponse(c, http.StatusInternalServerError, errResponseText)
	}

	err = json.NewDecoder(c.Request().Body).Decode(&header)
	if err != nil {
		log.Errorf("error decoding header request body: %v", err)
		return sendErrorResponse(c, http.StatusInternalServerError, errResponseText)
	}

	if app.User.ID == 0 || routeID == 0 {
		log.Errorf("error getting user id or route id:\n userID: %v \n routeID: %v")
		return sendErrorResponse(c, http.StatusBadRequest, errResponseText)
	}

	if c.Request().Method == "POST" {
		err = app.Db.SaveHeader(&header, app.User.ID, uint(routeID))
	}

	if c.Request().Method == "PUT" {
		err = app.Db.UpdateHeader(header, app.User.ID, uint(routeID))
	}

	if err != nil {
		log.Errorf("error updating header: %v", err)
		return sendErrorResponse(c, http.StatusBadRequest, errResponseText)
	}

	return sendOkResponse(c, header.ID)
}

func DeleteHeader(c echo.Context) error {
	app := c.(App)
	errResponseText := "couldn't delete header"

	routeID, err := strconv.Atoi(c.QueryParam("route"))
	if err != nil {
		log.Errorf("error converting route id query to int: %v", err)
		return sendErrorResponse(c, http.StatusInternalServerError, errResponseText)
	}
	headerName := c.QueryParam("name")

	if app.User.ID == 0 || routeID == 0 {
		log.Errorf("error getting user id or route id:\n userID: %v \n routeID: %v")
		return sendErrorResponse(c, http.StatusBadRequest, errResponseText)
	}

	err = app.Db.DeleteHeader(headerName, app.User.ID, uint(routeID))
	if err != nil {
		log.Errorf("error deleting header: %v", err)
		return sendErrorResponse(c, http.StatusBadRequest, errResponseText)
	}

	return sendOkResponse(c, "")
}
