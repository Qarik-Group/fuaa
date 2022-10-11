package core

import (
	uaa "github.com/cloudfoundry-community/go-uaa"
	"github.com/labstack/echo/v4"
)

type Services struct {
	Tokens TokenService
	Users  UserService
	Logger Logger
}

type Logger func(c echo.Context)

type TokenService interface {
	Create(uaa.User) (string, error)
	ByUser(uaa.User) (string, error)
	Exists(string) bool
	Reset()
}

type UserService interface {
	ByUsername(string) (uaa.User, error)
	Add(string, string) (uaa.User, error)
	Reset()
}

func (services *Services) Reset() {
	services.Tokens.Reset()
	services.Users.Reset()
}
