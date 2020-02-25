// Package blocker provide a list of predefined predicate methods
package blocker

import (
	"net/http"
	"strings"
)

// Method provides functionality to block a request based on a list of allowed HTTP methods.
// Empty allowedMethods will block
type Method struct {
	allowedMethods []string
}

// NewMethod is a constructor function for Method struct
func NewMethod(allowedMethods []string) *Method {
	return &Method{allowedMethods}
}

// Block is a predicate for method based request blocking.
func (m *Method) Block(r *http.Request) bool {
	block := true
	for _, v := range m.allowedMethods {
		if r.Method == strings.ToUpper(v) {
			block = false
			break
		}
	}
	return block
}
