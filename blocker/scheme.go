package blocker

import (
	"net/http"
	"strings"
)

// Scheme provides functionality to block a request based on a list of allowed HTTP schemes (http/https).
// Empty allowedSchemes will block
type Scheme struct {
	allowedSchemes []string
}

// NewScheme is a constructor function for Scheme struct
func NewScheme(allowedSchemes []string) *Scheme {
	return &Scheme{allowedSchemes}
}

// Block is a predicate for scheme based request blocking.
func (s *Scheme) Block(r *http.Request) bool {
	block := true
	for _, v := range s.allowedSchemes {
		if r.URL.Scheme == strings.ToLower(v) {
			block = false
			break
		}
	}
	return block
}
