package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/idanyas/go-pterodactyl"
)

// setup sets up a test HTTP server.
func setup() (mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	return mux, server.URL, server.Close
}

func testClient(t *testing.T, serverURL string) *pterodactyl.Client {
	t.Helper()
	c, err := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return c
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}
