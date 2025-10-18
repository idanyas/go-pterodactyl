package pterodactyl

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCheckResponse_Success(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Request:    &http.Request{},
	}

	err := CheckResponse(resp)
	if err != nil {
		t.Errorf("CheckResponse() returned error for 200 OK: %v", err)
	}
}

func TestCheckResponse_Error(t *testing.T) {
	body := `{
		"errors": [{
			"code": "ValidationException",
			"status": "400",
			"detail": "The email field is required.",
			"source": {
				"field": "email"
			}
		}]
	}`

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request: &http.Request{
			Method: "POST",
			URL:    mustParseURL("https://panel.example.com/api/users"),
		},
	}

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("CheckResponse() should return error for 400 status")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}

	if len(apiErr.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(apiErr.Errors))
	}

	if apiErr.Errors[0].Code != "ValidationException" {
		t.Errorf("Code = %s, want ValidationException", apiErr.Errors[0].Code)
	}

	if apiErr.Errors[0].Source == nil || apiErr.Errors[0].Source.Field != "email" {
		t.Error("expected source field 'email'")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "400") {
		t.Errorf("Error message should contain status code: %s", errMsg)
	}
	if !strings.Contains(errMsg, "email field is required") {
		t.Errorf("Error message should contain detail: %s", errMsg)
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		apiErr  *APIError
		wantMsg string
	}{
		{
			name: "with method and URL",
			apiErr: &APIError{
				StatusCode: 404,
				Method:     "GET",
				URL:        "https://panel.example.com/api/users/999",
				Errors: []ErrorDetail{
					{Code: "NotFound", Detail: "User not found"},
				},
			},
			wantMsg: "GET",
		},
		{
			name: "with field source",
			apiErr: &APIError{
				StatusCode: 422,
				Errors: []ErrorDetail{
					{
						Code:   "required",
						Detail: "The name field is required",
						Source: &ErrorSource{Field: "name"},
					},
				},
			},
			wantMsg: "field: name",
		},
		{
			name: "multiple errors",
			apiErr: &APIError{
				StatusCode: 422,
				Errors: []ErrorDetail{
					{Detail: "Error 1"},
					{Detail: "Error 2"},
				},
			},
			wantMsg: "Error 1; Error 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.apiErr.Error()
			if !strings.Contains(msg, tt.wantMsg) {
				t.Errorf("Error() = %q, want to contain %q", msg, tt.wantMsg)
			}
		})
	}
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
