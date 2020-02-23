![Build Status](https://github.com/psampaz/shield/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/psampaz/shield/branch/master/graph/badge.svg)](https://codecov.io/gh/psampaz/shield)
[![GoDoc](https://godoc.org/github.com/psampaz/shield?status.svg)](https://pkg.go.dev/github.com/psampaz/shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/psampaz/shield)](https://goreportcard.com/report/github.com/psampaz/shield)

# Shield

Shield is a net/http compatible middleware which blocks or allows request based on a predicate.
Shield replies back with a user defined response when the request is blocked.

# Usage

Below you can find a example of how to configure the Shield middleware in order to allow only requests with GET method, and reply back with 405 Method Not Allowed in any other case.

```go
package main

import (
	"net/http"

	"github.com/psampaz/shield"
)

func main() {

	shieldMiddleware := shield.New(shield.Options{
		Block: func(r *http.Request) bool {
			return r.Method != "GET"
		},
		Code:    http.StatusMethodNotAllowed,
		Headers: http.Header{"Content-Type": {"text/plain"}},
		Body:    []byte(http.StatusText(http.StatusMethodNotAllowed)),
	})
    
	helloWorldHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
    
	http.ListenAndServe(":8080", shieldMiddleware.Handler(helloWorldHandler))
}
```

```bash
$ curl -i -X GET localhost:8080

HTTP/1.1 200 OK
Date: Sat, 22 Feb 2020 10:03:35 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

hello world
```

```bash
$ curl -i -X POST localhost:8080

HTTP/1.1 400 Bad Request
Content-Type: text/plain
Date: Sat, 22 Feb 2020 10:02:31 GMT
Content-Length: 11

Bad Request
```

Passing a func as Block option, gives you access only in the current request. If there is a need to to use non request related data and functionality, you can you a stuct method with the same signature.

```go
package main

import (
	"net/http"

	"github.com/psampaz/shield"
)

type BlockLogic struct {
	ShouldBLock bool
}

func (b *BlockLogic) Block(r *http.Request) bool {
	return b.ShouldBLock
}

func main() {
	blockLogic := BlockLogic{true}
	shieldMiddleware := shield.New(shield.Options{
		Block:   blockLogic.Block,
		Code:    http.StatusBadRequest,
		Headers: http.Header{"Content-Type": {"text/plain"}},
		Body:    []byte(http.StatusText(http.StatusBadRequest)),
	})

	helloWorldHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	http.ListenAndServe(":8080", shieldMiddleware.Handler(helloWorldHandler))
}
```

## Options

Shield middleware can be configured with the following options:

```go
type Options struct {
	// Block is a predicate responsible for blocking the request.
	// Return true when the request should be blocked, false otherwise
	Block func(r *http.Request) bool
	// Code  is the status code of the response, sent when the request is blocked
	Code int
	// Headers are the headers of the response, sent when the request is blocked
	Headers http.Header
	// Body is the body of the response, sent when the request is blocked
	Body []byte
}
```

## Examples of block functions

### Block requests based on a HTTP header
```go
func(r *http.Request) bool {
    if r.Header.Get("X-Custom") != "" {
        return false
    }   

	return true
}
```
### Block requests based on HTTP method
```go
func(r *http.Request) bool {
    if r.Method == "GET" {
        return false
    }   

	return true
}
```
### Block requests based on HTTP scheme
```go
func(r *http.Request) bool {
    if r.URL.Sheme == "https" {
        return false
    }   

	return true
}
```
### Block requests based on query parameters
matched, err := regexp.MatchString(v, r.URL.Query().Get(k))
```go
func(r *http.Request) bool {
    // allow only request that have a query param named page which is a number
    matched, _ := regexp.MatchString(`\d+`, r.URL.Query().Get("page"))
    if matched {
        return false
    }
    return true
}
```

# Integration with popular routers

## Gorilla Mux

```go
package main

import (
	"net/http"

	"github.com/psampaz/shield"

	"github.com/gorilla/mux"
)

func main() {
	shieldMiddleware := shield.New(shield.Options{
		Block: func(r *http.Request) bool {
			return true
		},
		Code:    http.StatusMethodNotAllowed,
		Headers: http.Header{"Content-Type": {"text/plain"}},
		Body:    []byte(http.StatusText(http.StatusMethodNotAllowed)),
	})

	r := mux.NewRouter()
	r.Use(shieldMiddleware.Handler)
	r.HandleFunc("/", HelloHandler)

	http.ListenAndServe(":8080", r)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
```

## Chi

```go
package main

import (
	"net/http"

	"github.com/psampaz/shield"

	"github.com/go-chi/chi"
)

func main() {
	shieldMiddleware := shield.New(shield.Options{
		Block: func(r *http.Request) bool {
			return true
		},
		Code:    http.StatusMethodNotAllowed,
		Headers: http.Header{"Content-Type": {"text/plain"}},
		Body:    []byte(http.StatusText(http.StatusMethodNotAllowed)),
	})

	r := chi.NewRouter()
	r.Use(shieldMiddleware.Handler)
	r.Get("/", HelloHandler)

	http.ListenAndServe(":8080", r)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
```

## Echo

```go
package main

import (
	"net/http"

	"github.com/psampaz/shield"

	"github.com/labstack/echo"
)

func main() {
	shieldMiddleware := shield.New(shield.Options{
		Block: func(r *http.Request) bool {
			return true
		},
		Code:    http.StatusMethodNotAllowed,
		Headers: http.Header{"Content-Type": {"text/plain"}},
		Body:    []byte(http.StatusText(http.StatusMethodNotAllowed)),
	})

	e := echo.New()
	e.Use(echo.WrapMiddleware(shieldMiddleware.Handler))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.Start((":8080"))
}
```
