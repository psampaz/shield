// Package shield provides a net/http compatible middleware which blocks or allows request based on a predicate.
//
// Usage:
//	package main
//
//	import (
//		"net/http"
//
//		"github.com/psampaz/shield"
//	)
//
//	func main() {
//
//		shieldMiddleware := shield.New(shield.Options{
//			Block: func(r *http.Request) bool {
//				return r.Method != "GET"
//			},
//			Code:    http.StatusMethodNotAllowed,
//			Headers: http.Header{"Content-Type": {"text/plain"}},
//			Body:    []byte(http.StatusText(http.StatusMethodNotAllowed)),
//		})
//
//		helloWorldHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.Write([]byte("hello world"))
//		})
//
//		http.ListenAndServe(":8080", shieldMiddleware.Handler(helloWorldHandler))
//	}
package shield

import "net/http"

// Options holds configuration params for the shield
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

// Shield is an net/http compatible middleware
type Shield struct {
	options Options
}

// New creates a new Shield from Options
func New(o Options) *Shield {
	return &Shield{options: o}
}

// Handler middleware blocks a request based on a user defined predicate. When the request is blocked,
// the middleware responses back with user defined headers, status code and body. When the request is not blocked,
// the middleware just calls the next handler in the chain without altering the request.
func (s *Shield) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.options.Block(r) {
			for header, values := range s.options.Headers {
				for idx, value := range values {
					if idx == 0 {
						w.Header().Set(header, value)
					} else {
						w.Header().Add(header, value)
					}
				}
			}
			w.WriteHeader(s.options.Code)
			w.Write(s.options.Body)
			return
		}
		next.ServeHTTP(w, r)
	})
}
