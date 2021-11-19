package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func InitiateService(c echo.Context) error {
	app := c.(App)
	origin := c.Request().Header.Get("Origin")

	splitHost := strings.Split(c.Request().Host, ".")
	if splitHost[0] == "chapi" || len(splitHost) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	endpoint, err := app.Conn.GetRouteRequestData(splitHost[0], c.Request().URL.Path)
	if err != nil {
		log.Errorf("error getting project: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if len(endpoint.PermOrigins) == 0 || origin == os.Getenv("LOCAL_FRONTEND") {
		err = RunProxy(c, endpoint)
		return err
	}

	for _, endpointOrigin := range endpoint.PermOrigins {
		if endpointOrigin.URL == origin {
			err = RunProxy(c, endpoint)
			return err
		}
	}

	return nil
}

func RunProxy(c echo.Context, endpoint model.Endpoint) error {
	req, err := buildRequest(c.Request(), endpoint)
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
	resp.Body.Close()

	c.Response().Write(responseBody)

	return nil
}

func buildRequest(request *http.Request, endpoint model.Endpoint) (*http.Request, error) {
	urlQuery := request.URL.RawQuery

	destinationURL := fmt.Sprintf("%v%v?", endpoint.Destination, endpoint.Path)

	if urlQuery != "" || len(endpoint.Queries) != 0 {
		destinationURL += urlQuery

		for index, query := range endpoint.Queries {
			if index != 0 {
				destinationURL += "&"
			}

			if index == 0 && urlQuery != "" {
				destinationURL += "&"
			}

			destinationURL += fmt.Sprintf("%v=%v", query.Name, query.Value)
		}
	}

	req, err := http.NewRequest(endpoint.Method, destinationURL, nil)
	if err != nil {
		return nil, err
	}

	for _, header := range endpoint.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	return req, nil
}

func RunFrontendOrProxy(c echo.Context) error {
	host := c.Request().Host

	splitHost := strings.Split(host, ".")

	if splitHost[0] == "chapi" {
		HandleFrontend(c)
		return nil
	}

	InitiateService(c)

	return nil
}
