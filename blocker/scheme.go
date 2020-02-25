package blocker

import (
	"net/http"
	"strings"
)

// Scheme provides functionality to block a request based on a list of allowed HTTP schemes (http/https).
// Empty allowSchemes will block
type Scheme struct {
	allowSchemes []string
}

// NewScheme is a constructor function for Scheme struct
func NewScheme(allowSchemes []string) *Scheme {
	return &Scheme{allowSchemes}
}

// Block is a predicate for scheme based request blocking.
func (s *Scheme) Block(r *http.Request) bool {
	block := true
	for _, v := range s.allowSchemes {
		if r.URL.Scheme == strings.ToLower(v) {
			block = false
			break
		}
	}
	return block
}
