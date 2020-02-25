package blocker

import (
	"net/http"
	"regexp"
)

// Query provides functionality to block a request based on a query parameters.
// Empty queryParams will block
type Query struct {
	queryParams map[string]string
}

// NewQuery is a constructor function for Query struct
func NewQuery(queryParams map[string]string) *Query {
	return &Query{queryParams}
}

// Block is a predicate for query based request blocking.
func (m *Query) Block(r *http.Request) bool {
	if len(m.queryParams) == 0 {
		return true
	}

	block := false
	for k, v := range m.queryParams {
		matched, err := regexp.MatchString(v, r.URL.Query().Get(k))
		if err != nil || !matched {
			block = true
			break
		}
	}
	return block
}
