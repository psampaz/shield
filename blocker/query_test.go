package blocker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		qparams     map[string]string
		url         string
		shouldBlock bool
	}{
		{
			name:        "no query params",
			qparams:     nil,
			url:         "http://localhost/foo?q=1",
			shouldBlock: true,
		},
		{
			name:        "no query params",
			qparams:     map[string]string{},
			url:         "http://localhost/foo?q=1",
			shouldBlock: true,
		},
		{
			name:        "query param matches",
			qparams:     map[string]string{"q": "\\d"},
			url:         "http://localhost/foo?q=1",
			shouldBlock: false,
		},
		{
			name:        "query param matches",
			qparams:     map[string]string{"q": "value"},
			url:         "http://localhost/foo?q=value",
			shouldBlock: false,
		},
		{
			name:        "multiple query params match",
			qparams:     map[string]string{"q1": "value1", "q2": "value2"},
			url:         "http://localhost/foo?q1=value1&q2=value2",
			shouldBlock: false,
		},
		{
			name:        "query param does not match",
			qparams:     map[string]string{"q": `\d`},
			url:         "http://localhost/foo?q=a",
			shouldBlock: true,
		},
		{
			name:        "matcher without query params",
			qparams:     map[string]string{"q": `\d`},
			url:         "http://localhost/foo",
			shouldBlock: true,
		},
		{
			name:        "optional param",
			qparams:     map[string]string{"q": `^$|\d`},
			url:         "http://localhost/foo",
			shouldBlock: false,
		},
		{
			name:        "optional param",
			qparams:     map[string]string{"q": `^$|\d`},
			url:         "http://localhost/foo?q=",
			shouldBlock: false,
		},
		{
			name:        "optional param",
			qparams:     map[string]string{"q": `^$|\d`},
			url:         "http://localhost/foo?q=1",
			shouldBlock: false,
		},
		{
			name:        "multiple params",
			qparams:     map[string]string{"q1": `\d`, "q2": `\d`},
			url:         "http://localhost/foo?q1=1&q2=2",
			shouldBlock: false,
		},
		{
			name:        "multiple params, one optional",
			qparams:     map[string]string{"q1": `\d`, "q2": `^$|\d`},
			url:         "http://localhost/foo?q1=1",
			shouldBlock: false,
		},
		{
			name:        "multiple params, one optional",
			qparams:     map[string]string{"q1": `\d`, "q2": `^$|\d`},
			url:         "http://localhost/foo?q2=1",
			shouldBlock: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			blocker := NewQuery(tt.qparams)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			blocked := blocker.Block(req)

			if blocked != tt.shouldBlock {
				t.Errorf("Got status code %v, wanted %v", blocked, tt.shouldBlock)
			}
		})
	}
}

func ExampleQuery_Block() {
	m := NewQuery(map[string]string{"page": `\d+`})
	r := httptest.NewRequest(http.MethodPut, "http://localhost:8080?page=a", nil)
	blocked := m.Block(r)
	fmt.Printf("%t", blocked)
	// Output: true
}
