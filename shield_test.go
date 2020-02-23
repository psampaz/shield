package shield

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var shieldBody []byte = []byte("shield")
var shieldCode int = http.StatusBadRequest

// shieldResponseHeaders are the headers of the response when the request is blocked
var shieldResponseHeaders http.Header = map[string][]string{
	"Shield-H1": {
		"shield-h1-v1", "shield-h1-v2",
	},
	"Shield-H2": {
		"shield-h2-v1", "shield-h2-v2",
	},
	"Shield-H3": {
		"shield-h3-v1",
	},
}

var nextBody []byte = []byte("next")
var nextCode int = http.StatusOK

// nextResponseHeaders are the headers of the response of the next handler(when the request is not blocked)
var nextResponseHeaders http.Header = map[string][]string{
	"Next-H1": {
		"next-h1-v1", "next-h1-v2",
	},
	"Next-H2": {
		"next-h2-v1", "next-h2-v2",
	},
	"Next-H3": {
		"next-h3-v1",
	},
}

// this is the next handler that is wrapped by the shield handler
var nextHandlerFunc http.HandlerFunc = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	for header, values := range nextResponseHeaders {
		for idx, value := range values {
			if idx == 0 {
				w.Header().Set(header, value)
			} else {
				w.Header().Add(header, value)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nextBody)
})

func TestBlock(t *testing.T) {
	tests := []struct {
		name        string
		description string
		options     Options
		next        http.HandlerFunc
		request     *http.Request
		wantCode    int
		wantHeaders http.Header
		wantBody    []byte
	}{
		{
			name:        "Block request",
			description: "Response code, headers and body should be the ones configured in the shield",
			options: Options{
				Block:   func(r *http.Request) bool { return true },
				Code:    shieldCode,
				Headers: shieldResponseHeaders,
				Body:    shieldBody,
			},
			next:        nextHandlerFunc,
			wantCode:    shieldCode,
			wantHeaders: shieldResponseHeaders,
			wantBody:    shieldBody,
		},
		{
			name:        "Allow request",
			description: "Next response code, headers and body should not be altered",
			options: Options{
				Block:   func(r *http.Request) bool { return false },
				Code:    shieldCode,
				Headers: shieldResponseHeaders,
				Body:    shieldBody,
			},
			next:        nextHandlerFunc,
			wantCode:    nextCode,
			wantHeaders: nextResponseHeaders,
			wantBody:    nextBody,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shield := New(tt.options)

			w := httptest.NewRecorder()
			shield.Handler(tt.next).ServeHTTP(w, httptest.NewRequest("GET", "http://localhost", nil))

			resp := w.Result()

			gotCode := resp.StatusCode
			if gotCode != tt.wantCode {
				t.Errorf("Got status code %d, wanted %d", gotCode, tt.wantCode)
			}

			gotHeader := resp.Header
			if !reflect.DeepEqual(tt.wantHeaders, gotHeader) {
				t.Errorf("Got headers %v, wanted %v", gotHeader, tt.wantHeaders)
			}

			gotBody, _ := ioutil.ReadAll(resp.Body)
			if !bytes.Equal(tt.wantBody, gotBody) {
				t.Errorf("Got body %q, wanted %q", string(gotBody), string(tt.wantBody))
			}
		})
	}
}
