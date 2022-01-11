package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

var (
	ENVIRONMENT     = os.Getenv("ENV")
	protectedRoutes = []string{
		"localhost",
		"www",
		"chapi",
		"webmail",
	}
)

func RunFrontendOrProxy(c echo.Context) error {
	hostDomain := getSubdomain(c.Request().Host)

	if hostDomain != "" && isDomainProtected(hostDomain) {
		ServeStaticAssets((c))
		return nil
	}

	method := c.Request().Method

	if method != http.MethodPost && method != http.MethodGet {
		return nil
	}

	err := InitiateService(c, hostDomain)
	if err != nil {
		return err
	}

	return nil
}

func ServeStaticAssets(c echo.Context) error {
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

	c.Response().Header().Add("Cache-Control", "max-age=86400")
	c.Response().Writer.Write(fileData)
	return nil
}

func InitiateService(c echo.Context, domain string) error {
	app := c.(App)
	origin := c.Request().Header.Get("Origin")
	fmt.Printf("Initiating service: %v from %v", domain, origin)

	endpoint, err := app.Conn.GetRouteRequestData(domain, c.Request().URL.Path)
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

	responseBody := bytes.NewBuffer(nil)
	responseSize, err := io.Copy(responseBody, resp.Body)
	if err != nil {
		log.Errorf("error reading response body: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if responseSize > 30000000 {
		c.JSON(http.StatusOK, "Response too big")
		return nil
	}

	c.Response().Write(responseBody.Bytes())

	return nil
}

func joinRequestBodies(request *http.Request, endpoint model.Endpoint) (*bytes.Reader, error) {
	requestBody := make(map[string]interface{})
	endpointBody := make(map[string]interface{})

	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil && strings.ToLower(err.Error()) != "eof" {
		log.Warnf("error reading request body: %v", err)
	}

	err = json.NewDecoder(bytes.NewBufferString(endpoint.Body)).Decode(&endpointBody)
	if err != nil && strings.ToLower(err.Error()) != "eof" {
		log.Warnf("error reading endpoint body: %v", err)
	}

	for k, v := range endpointBody {
		requestBody[k] = v
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	return reader, nil
}

func joinRequestQueries(request *http.Request, endpoint model.Endpoint) (string, error) {
	urlQuery := request.URL.RawQuery

	fullURL := endpoint.Destination + "?"

	if urlQuery != "" || len(endpoint.Queries) != 0 {
		fullURL += urlQuery

		for index, query := range endpoint.Queries {
			if index != 0 {
				fullURL += "&"
			}

			if index == 0 && urlQuery != "" {
				fullURL += "&"
			}

			fullURL += fmt.Sprintf("%v=%v", query.Name, query.Value)
		}
	}

	return fullURL, nil
}

func buildRequest(request *http.Request, endpoint model.Endpoint) (*http.Request, error) {
	body, err := joinRequestBodies(request, endpoint)
	if err != nil {
		return nil, err
	}

	destinationURL, err := joinRequestQueries(request, endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(endpoint.Method, destinationURL, body)
	if err != nil {
		return nil, err
	}

	for _, header := range endpoint.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	return req, nil
}

func getSubdomain(url string) string {
	var domain string

	if strings.Contains(url, "http://") {
		url = url[7:]
	}
	if strings.Contains(url, "https://") {
		url = url[8:]
	}

	domain = strings.Split(url, ".")[0]

	return domain
}

func isDomainProtected(domain string) bool {
	for _, route := range protectedRoutes {
		if strings.Contains(domain, route) {
			return true
		}
	}
	return false
}
