package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/starkandwayne/fuaa/core"
	"github.com/starkandwayne/fuaa/http/routes/registry"
)

type GetLogin struct {
	path     string
	verb     registry.Verb
	ready    bool
	services *core.Services
	urls     map[string]string
}

func NewGetLogin(services *core.Services, urls map[string]string) *GetLogin {
	return &GetLogin{
		path:     "/login",
		verb:     registry.GET,
		ready:    true,
		services: services,
		urls:     urls,
	}
}

func (route *GetLogin) Register(router *echo.Echo) {
	registry.Register(router, route)
}

func (route *GetLogin) Path() string {
	return route.path
}

func (route *GetLogin) Ready() bool {
	return route.ready
}

func (route *GetLogin) Verb() registry.Verb {
	return route.verb
}

func (route *GetLogin) Handle(c echo.Context) error {
	route.services.Logger(c)

	uaaURL, ok := route.urls["uaaURL"]
	if !ok {
		uaaURL = "https://uaa.example.com"
	}

	template := `{
  "app": {
    "version": "4.31.0-SNAPSHOT"
  },
  "links": {
    "uaa": "%s",
    "passwd": "/forgot_password",
    "login": "%s",
    "register": "/create_account"
  },
  "zone_name": "uaa",
  "entityID": "%s",
  "commit_id": "2e27fa7",
  "idpDefinitions": {},
  "prompts": {
    "username": [
      "text",
      "Email"
    ],
    "password": [
      "password",
      "Password"
    ]
  },
  "timestamp": "%s"
}`

	now := time.Now().UTC().String()
	output := fmt.Sprintf(template, uaaURL, uaaURL, uaaURL, now)
	return c.JSONBlob(
		http.StatusOK,
		[]byte(output),
	)
}
