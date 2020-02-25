// Package blocker provide a list of predefined predicate methods
package blocker

import "net/http"

// Method provides functionality to block a request based on a list of allowed HTTP methods.
// Empty allowMethods will block
type Method struct {
	allowMethods []string
}

// NewMethod is a constructor function for Method struct
func NewMethod(allowMethods []string) *Method {
	return &Method{allowMethods}
}

// Block is a predicate for method based request blocking.
func (m *Method) Block(r *http.Request) bool {
	block := true
	for _, v := range m.allowMethods {
		if r.Method == v {
			block = false
			break
		}
	}
	return block
}
