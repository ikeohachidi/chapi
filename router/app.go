package router

import (
	"embed"
	"net/http"
	"os"

	goSession "github.com/gorilla/sessions"
	"github.com/ikeohachidi/chapi/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var store = goSession.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type App struct {
	echo.Context
	Conn model.Conn
	User model.User
	Fs   embed.FS
}

type Response struct {
	Data       interface{} `json:"data"`
	Successful bool        `json:"successful"`
}

func errRedirect(c echo.Context, url string, err error) {
	log.Error(err)
	c.Redirect(http.StatusInternalServerError, url)
}

func sendErrorResponse(c echo.Context, statusCode int, responseText string) error {
	return c.JSON(statusCode, Response{
		Data:       responseText,
		Successful: false,
	})
}

func sendOkResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Data:       data,
		Successful: true,
	})
}
