package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/ikeohachidi/chapi-be/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func HandleFrontend(c echo.Context) error {
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

func StartProxy(c echo.Context) error {
	// blog.chapi.com/openmap
	app := c.(App)

	splitHost := strings.Split(c.Request().Host, ".")
	if splitHost[0] == "chapi" || len(splitHost) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	route, err := app.Db.GetRouteFromNameAndPath(splitHost[0], c.Request().URL.Path)
	if err != nil {
		log.Errorf("error getting project: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	destinationURL := buildURL(c.Request(), route)

	req, err := http.NewRequest(route.Type, destinationURL, nil)
	if err != nil {
		log.Errorf("error creating new request: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("error making http request with default client: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	c.Response().Write(responseBody)

	return nil
}

func buildURL(request *http.Request, route model.Route) string {
	urlQuery := request.URL.RawQuery

	destinationURL := fmt.Sprintf("%v%v?", route.Destination, route.Path)

	if urlQuery != "" || len(route.Queries) != 0 {
		destinationURL += urlQuery

		for index, query := range route.Queries {
			if index != 0 {
				destinationURL += "&"
			}

			if index == 0 && urlQuery != "" {
				destinationURL += "&"
			}

			destinationURL += fmt.Sprintf("%v=%v", query.Name, query.Value)
		}
	}

	return destinationURL
}

func RunFrontendOrProxy(c echo.Context) error {
	host := c.Request().Host

	splitHost := strings.Split(host, ".")

	if splitHost[0] == "chapi" {
		HandleFrontend(c)
		return nil
	}

	StartProxy(c)

	return nil
}
