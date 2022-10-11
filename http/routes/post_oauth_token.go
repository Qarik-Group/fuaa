package routes

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"codeberg.org/ess/fuaa/core"
	"codeberg.org/ess/fuaa/http/routes/registry"
)

type PostOauthToken struct {
	path     string
	verb     registry.Verb
	ready    bool
	services *core.Services
	urls     map[string]string
}

func NewPostOauthToken(services *core.Services, urls map[string]string) *PostOauthToken {
	return &PostOauthToken{
		path:     "/oauth/token",
		verb:     registry.POST,
		ready:    true,
		services: services,
		urls:     urls,
	}
}

func (route *PostOauthToken) Register(router *echo.Echo) {
	registry.Register(router, route)
}

func (route *PostOauthToken) Path() string {
	return route.path
}

func (route *PostOauthToken) Ready() bool {
	return route.ready
}

func (route *PostOauthToken) Verb() registry.Verb {
	return route.verb
}

func (route *PostOauthToken) Handle(c echo.Context) error {
	route.services.Logger(c)

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			"body read fail",
		)
	}

	params, err := url.ParseQuery(string(body))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			"couldn't parse the payload",
		)
	}

	grantType := params.Get("grant_type")

	switch grantType {
	case "password":
		username := params.Get("username")
		user, err := route.services.Users.ByUsername(username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		password := params.Get("password")
		if user.Password != password {
			return c.JSON(http.StatusUnauthorized, nil)
		}

		token, err := route.services.Tokens.ByUser(user)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, nil)
		}

		return c.JSON(
			http.StatusOK,
			map[string]string{
				"token_type":    "bearer",
				"access_token":  string(token),
				"refresh_token": string(token),
			},
		)
	case "refresh_token":
		token := params.Get("refresh_token")
		if len(token) == 0 {
			return c.JSON(http.StatusInternalServerError, "no refresh_token given")
		}

		if !route.services.Tokens.Exists(token) {
			return c.JSON(http.StatusUnauthorized, "no such token")
		}

		return c.JSON(
			http.StatusOK,
			map[string]string{
				"token_type":    "bearer",
				"access_token":  token,
				"refresh_token": token,
			},
		)
	default:
		return c.JSON(http.StatusInternalServerError, "unimplemented grant type")
	}
}
