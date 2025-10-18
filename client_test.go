package pterodactyl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// setup sets up a test HTTP server along with a Client that is
// configured to talk to that server. It returns the mux, server URL, and teardown function.
func setup() (mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	return mux, server.URL, server.Close
}

func testClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client, err := New(serverURL, WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return client
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

func TestNew(t *testing.T) {
	c, err := New("https://panel.example.com", WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("New() returned an unexpected error: %v", err)
	}

	wantBaseURL, _ := url.Parse("https://panel.example.com/api/")
	if !reflect.DeepEqual(c.baseURL, wantBaseURL) {
		t.Errorf("New() baseURL is %v, want %v", c.baseURL, wantBaseURL)
	}

	if c.apiKey != "test-key" {
		t.Errorf("New() apiKey is %s, want %s", c.apiKey, "test-key")
	}

	if c.Application() == nil {
		t.Error("New() Application client is nil")
	}

	if c.Client() == nil {
		t.Error("New() Client client is nil")
	}
}

func TestNew_emptyURL(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Error("New() with empty URL should return an error")
	}
}

func TestClient_Do(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	client := testClient(t, serverURL)

	type foo struct {
		A string
	}

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", "Application/vnd.pterodactyl.v1+json")
		testHeader(t, r, "Authorization", "Bearer test-key")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	body := new(foo)
	_, err := client.Do(context.Background(), http.MethodGet, "", nil, body)
	if err != nil {
		t.Fatalf("Do() returned error: %v", err)
	}

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestClient_Do_httpError(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	client := testClient(t, serverURL)

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	_, err := client.Do(context.Background(), http.MethodGet, "", nil, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got nil.")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
}

func TestClient_Do_withBody(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	client := testClient(t, serverURL)

	type RequestBody struct {
		Name string `json:"name"`
	}
	type ResponseBody struct {
		Success bool `json:"success"`
	}

	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")

		var reqBody RequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("decoding request body failed: %v", err)
		}
		if reqBody.Name != "test-name" {
			t.Errorf("request body name is %s, want %s", reqBody.Name, "test-name")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ResponseBody{Success: true})
	})

	reqBody := RequestBody{Name: "test-name"}
	var respBody ResponseBody
	resp, err := client.Do(context.Background(), http.MethodPost, "test", &reqBody, &respBody)
	if err != nil {
		t.Fatalf("Do() returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("response status code is %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if !respBody.Success {
		t.Error("response body success is false, want true")
	}
}