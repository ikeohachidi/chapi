package main

import (
	"net/http"
	"os"

	goSession "github.com/gorilla/sessions"
	"github.com/ikeohachidi/chapi-be/model"
	"github.com/ikeohachidi/chapi-be/router"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()
	db := model.Connect()

	store := goSession.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := store.Get(c.Request(), "chapi_session")
			var user model.User

			if err != nil {
				log.Fatal("couldn't get chapi session")
				return c.JSON(http.StatusBadRequest, router.Response{"Couldn't get chapi session", false})
			}

			if _, ok := session.Values["access_token"]; ok {
				user.Email = session.Values["email"].(string)
				user.ID = session.Values["id"].(uint)
			}

			cc := router.App{
				c,
				db,
				user,
			}

			return next(cc)
		}
	})

	e.GET("/auth/github", router.OauthGithub)
	e.GET("/auth/github/redirect", router.OauthGithubRedirect)
	e.GET("/auth/logout", router.Logout)

	e.GET("/project", router.GetUserProjects)
	e.GET("/project/all", router.ListProjects)
	e.POST("/project", router.CreateProject)
	e.DELETE("/project/:id", router.DeleteProject)

	e.GET("/route/:projectID", router.GetProjectRoutes)
	e.POST("/route", router.SaveRoute)
	e.DELETE("/route", router.DeleteRoute)

	e.GET("/query/:routeID", router.GetRouteQueries)
	e.POST("/query", router.SaveQuery)
	e.DELETE("/query/:id", router.DeleteRouteQuery)

	e.Logger.Fatal(e.Start(":1333"))
}
