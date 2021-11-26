package router

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/ikeohachidi/chapi/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	FRONTEND_URL = os.Getenv("LOCAL_FRONTEND")
	SERVER_URL   = os.Getenv("LOCAL_SERVER")
)

var githubConfig = oauth2.Config{
	ClientID:     os.Getenv("CHAPI_GITHUB_ID"),
	ClientSecret: os.Getenv("CHAPI_GITHUB_SECRET"),
	Scopes:       []string{"user:email"},
	RedirectURL:  SERVER_URL + "/auth/github/redirect",
	Endpoint:     github.Endpoint,
}

func createState() string {
	//TODO: change this is for dev purpose only
	return "hello"
}

// OauthGithub will Redirect to the Github Authorization Page
func OauthGithub(c echo.Context) error {
	state := createState()

	authCode := githubConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	http.Redirect(c.Response().Writer, c.Request(), authCode, http.StatusTemporaryRedirect)

	return nil
}

// OauthGithubRedirect handles redirect requests from the Github Authorization page
func OauthGithubRedirect(c echo.Context) error {
	cc := c.(App)

	code := c.FormValue("code")

	// Exchange the user information for an access_token
	token, err := githubConfig.Exchange(context.Background(), code)
	errRedirect(c, SERVER_URL, err)

	// Get the user information from the github api
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	errRedirect(c, FRONTEND_URL, err)

	req.Header.Set("Authorization", "token "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	errRedirect(c, SERVER_URL, err)
	defer res.Body.Close()

	var user model.User

	err = json.NewDecoder(res.Body).Decode(&user)
	errRedirect(c, SERVER_URL, err)

	err = user.Create(cc.Conn.Db)
	errRedirect(c, SERVER_URL, err)

	// set cookie
	session, _ := store.Get(c.Request(), "chapi_session")
	session.Values["id"] = user.ID
	session.Values["email"] = user.Email
	session.Values["access_token"] = token.AccessToken
	session.Save(c.Request(), c.Response().Writer)

	c.Redirect(http.StatusMovedPermanently, SERVER_URL)

	return nil
}

// Logout deletes a users session
func Logout(c echo.Context) error {
	session, err := store.Get(c.Request(), "chapi_session")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{"Couldn't logout user", false})
	}

	delete(session.Values, "id")
	delete(session.Values, "email")
	delete(session.Values, "access_token")

	deleteCookie := http.Cookie{
		Name:    "chapi_session",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	}

	http.SetCookie(c.Response().Writer, &deleteCookie)
	return c.JSON(http.StatusOK, Response{"Logout successful", true})
}

func GetAuthenticatedUser(c echo.Context) (err error) {
	session, _ := store.Get(c.Request(), "chapi_session")

	if _, ok := session.Values["id"]; !ok {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, Response{
		Data: model.User{
			ID:    session.Values["id"].(uint),
			Email: session.Values["email"].(string),
		},
		Successful: true,
	})
}
