package application_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setup sets up a test HTTP server.
func setup() (mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	return mux, server.URL, server.Close
}

// testMethod checks the HTTP method of a request.
func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}
