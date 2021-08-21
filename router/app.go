package router

import (
	"net/http"
	"os"

	goSession "github.com/gorilla/sessions"
	"github.com/ikeohachidi/chapi-be/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var store = goSession.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type App struct {
	echo.Context
	Db model.Conn
}

type Response struct {
	Data       interface{} `json:"message"`
	Successful bool        `json:"successful"`
}

func errRedirect(c echo.Context, url string, err error) {
	log.Error(err)
	c.Redirect(http.StatusInternalServerError, url)
}
