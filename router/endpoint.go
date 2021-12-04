package router

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ikeohachidi/chapi/lib"
	"github.com/ikeohachidi/chapi/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func HandleFrontend(c echo.Context) error {
	app := c.(App)

	url := c.Request().URL.String()

	var file fs.File
	var fileData []byte
	var err error

	isFrontendPage := strings.Contains(url, "dashboard") && filepath.Ext(url) == ""

	if url == "/" || isFrontendPage {
		file, err = app.Fs.Open("frontend/dist/index.html")
		fileData, _ = ioutil.ReadAll(file)
	} else {
		c.Response().Header().Add("Content-Type", lib.DetectContentType(url))

		file, err = app.Fs.Open("frontend/dist" + url)
		fileData, _ = ioutil.ReadAll(file)
	}

	if err != nil {
		log.Errorf("couldn't parse %v file: %v", url, err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	c.Response().Writer.Write(fileData)
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

	destinationURL := endpoint.Destination + "?"

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

	if strings.Contains(CHAPI_SERVER_URL, splitHost[0]) {
		HandleFrontend(c)
		return nil
	}

	InitiateService(c)

	return nil
}
