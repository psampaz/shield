package blocker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethod(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		allowMethods  []string
		requestMethod string
		shouldBlock   bool
	}{
		{
			name:          "no methods to match",
			allowMethods:  nil,
			requestMethod: http.MethodGet,
			shouldBlock:   true,
		},
		{
			name:          "no methods to match",
			allowMethods:  []string{},
			requestMethod: http.MethodGet,
			shouldBlock:   true,
		},
		{
			name:          "single requestMethod match",
			allowMethods:  []string{http.MethodGet},
			requestMethod: http.MethodGet,
			shouldBlock:   false,
		},
		{
			name:          "single requestMethod mismatch",
			allowMethods:  []string{http.MethodPost},
			requestMethod: http.MethodGet,
			shouldBlock:   true,
		},
		{
			name:          "multiple requestMethod match",
			allowMethods:  []string{http.MethodGet, http.MethodPost},
			requestMethod: http.MethodGet,
			shouldBlock:   false,
		},
		{
			name:          "multiple requestMethod mismatch",
			allowMethods:  []string{http.MethodGet, http.MethodPost},
			requestMethod: http.MethodPut,
			shouldBlock:   true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			blocker := NewMethod(tt.allowMethods)

			req := httptest.NewRequest(tt.requestMethod, "http://localhost/", nil)
			blocked := blocker.Block(req)

			if blocked != tt.shouldBlock {
				t.Errorf("Got status code %v, wanted %v", blocked, tt.shouldBlock)
			}
		})
	}
}

func ExampleMethod_Block() {
	m := NewMethod([]string{http.MethodGet, http.MethodPost})
	r := httptest.NewRequest(http.MethodPut, "http://localhost:8080/", nil)
	blocked := m.Block(r)
	fmt.Printf("%t", blocked)
	// Output: true
}
