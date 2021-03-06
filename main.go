package main

import (
	"embed"
	"os"

	goSession "github.com/gorilla/sessions"
	"github.com/ikeohachidi/chapi/model"
	"github.com/ikeohachidi/chapi/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var (
	PORT = os.Getenv("PORT")

	//go:embed frontend/dist
	FS embed.FS
)

func main() {
	e := echo.New()
	db := model.Connect()

	store := goSession.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := store.Get(c.Request(), "chapi_session")
			var user model.User

			if err != nil {
				log.Errorf("error getting chapi session: %v", err)
			}

			if _, ok := session.Values["access_token"]; ok {
				user.Email = session.Values["email"].(string)
				user.ID = session.Values["id"].(uint)
			}

			cc := router.App{
				Context: c,
				Conn:    db,
				User:    user,
				Fs:      FS,
			}

			return next(cc)
		}
	})

	e.Any("/*", router.RunFrontendOrProxy)

	e.GET("/auth/github", router.OauthGithub)
	e.GET("/auth/github/redirect", router.OauthGithubRedirect)
	e.GET("/auth/logout", router.Logout)
	e.GET("/auth/user", router.GetAuthenticatedUser)

	e.GET("/project", router.GetUserProjects)
	e.GET("/project/all", router.ListProjects)
	e.POST("/project/exists", router.DoesProjectExist)
	e.POST("/project", router.SaveProject)
	e.PUT("/project", router.SaveProject)
	e.DELETE("/project/:id", router.DeleteProject)

	e.GET("/route", router.GetProjectRoutes)
	e.POST("/route", router.SaveRoute)
	e.PUT("/route", router.SaveRoute)
	e.DELETE("/route", router.DeleteRoute)

	e.GET("/query", router.GetQueries)
	e.POST("/query", router.SaveQuery)
	e.PUT("/query", router.SaveQuery)
	e.DELETE("/query", router.DeleteQuery)

	e.GET("/header", router.GetHeaders)
	e.POST("/header", router.SaveHeader)
	e.PUT("/header", router.SaveHeader)
	e.DELETE("/header", router.DeleteHeader)

	e.GET("/perm_origin", router.GetPermOrigins)
	e.POST("/perm_origin", router.SavePermOrigins)
	e.PUT("/perm_origin", router.SavePermOrigins)
	e.DELETE("/perm_origin", router.DeletePermOrigin)

	e.GET("/merge_options", router.GetMergeOptions)
	e.PUT("/merge_options", router.SaveMergeOptions)

	e.Logger.Fatal(e.Start(":" + PORT))
}
