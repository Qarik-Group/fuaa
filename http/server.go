package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codeberg.org/ess/fuaa/core"
	"codeberg.org/ess/fuaa/http/routes"
)

func Server(bind string, services *core.Services, urls map[string]string) *http.Server {
	e := echo.New()

	routes.Register(e, services, urls)

	e.GET("*", func(c echo.Context) error {
		fmt.Println("Got a GET request I didn't recognize:", c.Request().URL)
		services.Logger(c)

		return c.JSON(http.StatusInternalServerError, "")
	})

	e.POST("*", func(c echo.Context) error {
		fmt.Println("Got a POST request I didn't recognize:", c.Request().URL)
		services.Logger(c)

		return c.JSON(http.StatusInternalServerError, "")
	})

	server := e.Server
	server.Addr = bind

	return server
}
