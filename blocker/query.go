package blocker

import (
	"net/http"
	"regexp"
)

// Query provides functionality to block a request based on a query parameters.
// Empty qparams will block
type Query struct {
	qparams map[string]string
}

// NewQuery is a constructor function for Query struct
func NewQuery(qparams map[string]string) *Query {
	return &Query{qparams}
}

// Block is a predicate for query based request blocking.
func (m *Query) Block(r *http.Request) bool {
	if len(m.qparams) == 0 {
		return true
	}

	block := false
	for k, v := range m.qparams {
		matched, err := regexp.MatchString(v, r.URL.Query().Get(k))
		if err != nil || !matched {
			block = true
			break
		}
	}
	return block
}
