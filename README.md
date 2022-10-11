# fuaa #

[![Go Report Card](https://goreportcard.com/badge/codeberg.org/ess/fuaa)](https://goreportcard.com/report/codeberg.org/ess/fuaa)
[![Documentation](https://godoc.org/codeberg.org/ess/fuaa?status.svg)](http://godoc.org/codeberg.org/ess/fuaa)

A faux UAA server that acts just enough like a real UAA server to be dangerous.

## Installation ##

You can use fuaa either as a standalone service that runs on 0.0.0.0:8001, or you can use it programmatically in your test suite. To install the standalone executable, do this:

```
go get -u codeberg.org/ess/fuaa/cmd/fuaa
```

For programmatic usage, treat it like you would any library.

## Basics ##

At present, fuaa only responds to the following routes:

* `GET /login` - Get metadata for the service itself (used by the `cf` CLI, primarily)
* `POST /oauth/token` - Handle password logins and token refreshes

### Programmatic Usage ###

Let's say that you're writing an app that needs to interact with UAA in production, and you want to use fuaa to help builod out those interactions in development, AND you're writing that app in Go.

My friend, the fuaa library is for you.

You need a few things to make this happen:

* a `Services` collection
* a `map[string]string` of URLs (most specifically, `"uaaURL" : "http://localhost:8001"` or some such)
* a bind string (`hostname:port`) for the server to listen on

The last two of those items should coincide with each other.

```go
import (
  "testing"

  "codeberg.org/ess/fuaa/core"
  "codeberg.org/ess/fuaa/memory"
  "codeberg.org/ess/fuaa/http"
)

var services *core.Services

func MyTestSetup() {
  services = memory.NewServices()
  urls := map[string]string{"uaaURL" : "http://localhost:8675"}
  server := http.Server("0.0.0.0:8675", services, urls)

  go server.ListenAndServe()
}

// Run this between your test cases to clear the server's records
func ResetMocks() {
  services.Reset()
}
```

You should now be able to point your app at port 8675 on your local machine and, with any luck, be able to basically use it as you would UAA.

### Standalone Usage ###

Let's say that you're writing an app that needs to interact with UAA in production, and you want to use fuaa to help build out those interactions in development, but you're not writing that app in Go.

My friend, the standalone server is for you, and doubly so if you're using Docker Compose or something similar to manage the external dependencies for your development environment.

The only requirement for execution is to set the `UAA_ENDPOINT` environment variable, and it should be the URL that you want fuaa to advertise as itself. For example, when testing the routes out locally, I ran this to start the service:

```
UAA_ENDPOINT=http://localhost:8001 fuaa
```

Unfortunately, there is not presently a way to reset the standalone server's memory without fully restarting the server, so you'll need to keep that in mind for your test cases.

## Trivia ##

* This library was extracted from another project when we realized that we needed it for more than one project.
* It began its life as jeuaa: Just Enough (TM) UAA. The name was changed just because "fuaa" is *way* easier to type.

