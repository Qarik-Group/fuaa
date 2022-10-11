package registry

import (
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	locker        sync.Mutex
	ready         = make([]*routespec, 0)
	unimplemented = make([]*routespec, 0)
)

type routespec struct {
	Path string `json:"path"`
	Verb string `json:"verb"`
}

func regReady(path string, verb string) {
	ready = append(ready, &routespec{Path: path, Verb: verb})
}

func regUnimplemented(path string, verb string) {
	unimplemented = append(unimplemented, &routespec{Path: path, Verb: verb})
}

type Verb string

const (
	GET    Verb = "GET"
	POST   Verb = "POST"
	PUT    Verb = "PUT"
	PATCH  Verb = "PATCH"
	DELETE Verb = "DELETE"
	HEAD   Verb = "HEAD"
)

type Route interface {
	Handle(echo.Context) error
	Ready() bool
	Path() string
	Verb() Verb
}

type regfunc func(*echo.Echo, Route)

func Register(router *echo.Echo, route Route) {
	// While we shouldn't ever actually have a race condition, I'd rather be
	// safe than sorry.
	locker.Lock()
	defer locker.Unlock()

	if route.Ready() {
		// Actually tie implemented routes to the router
		router.Add(string(route.Verb()), route.Path(), route.Handle)

		// Register the route as ready
		regReady(route.Path(), string(route.Verb()))
	} else {
		// Register the route as unimplemented
		regUnimplemented(route.Path(), string(route.Verb()))
	}
}

func Registered() map[string][]*routespec {
	locker.Lock()
	defer locker.Unlock()

	output := make(map[string][]*routespec)

	output["ready"] = ready
	output["unimplemented"] = unimplemented

	return output
}
