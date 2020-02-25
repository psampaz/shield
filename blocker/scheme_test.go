package blocker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScheme(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		allowSchemes  []string
		requestScheme string
		shouldBlock   bool
	}{
		{
			name:          "no allowed schemes",
			allowSchemes:  nil,
			requestScheme: "http",
			shouldBlock:   true,
		},
		{
			name:          "no allowed schemes",
			allowSchemes:  []string{},
			requestScheme: "http",
			shouldBlock:   true,
		},
		{
			name:          "allow https, request is http",
			allowSchemes:  []string{"https"},
			requestScheme: "http",
			shouldBlock:   true,
		},
		{
			name:          "allow HTTPS, request is http",
			allowSchemes:  []string{"HTTPS"},
			requestScheme: "http",
			shouldBlock:   true,
		},
		{
			name:          "allow http, request is http",
			allowSchemes:  []string{"http"},
			requestScheme: "http",
			shouldBlock:   false,
		},
		{
			name:          "allow https, request is https",
			allowSchemes:  []string{"https"},
			requestScheme: "https",
			shouldBlock:   false,
		},
		{
			name:          "allow http, request is https",
			allowSchemes:  []string{"http"},
			requestScheme: "https",
			shouldBlock:   true,
		},
		{
			name:          "allow HTTP, request is https",
			allowSchemes:  []string{"HTTP"},
			requestScheme: "https",
			shouldBlock:   true,
		},
		{
			name:          "multiple schemes",
			allowSchemes:  []string{"http", "https"},
			requestScheme: "https",
			shouldBlock:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			blocker := NewScheme(tt.allowSchemes)

			req := httptest.NewRequest(http.MethodGet, tt.requestScheme+"://localhost/", nil)
			blocked := blocker.Block(req)

			if blocked != tt.shouldBlock {
				t.Errorf("Got status code %v, wanted %v", blocked, tt.shouldBlock)
			}
		})
	}
}

func ExampleScheme_Block() {
	m := NewScheme([]string{"https"})
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	blocked := m.Block(r)
	fmt.Printf("%t", blocked)
	// Output: true
}
