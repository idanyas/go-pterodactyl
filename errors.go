package pterodactyl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/idanyas/go-pterodactyl/transport"
)

// ErrorSource provides an optional object pointing to the specific field that caused the error.
type ErrorSource struct {
	Field string `json:"field"`
}

// ErrorDetail represents a single error object from the Pterodactyl API.
type ErrorDetail struct {
	Code   string       `json:"code"`
	Status string       `json:"status"`
	Detail string       `json:"detail"`
	Source *ErrorSource `json:"source,omitempty"`
}

// errorResponse is the structure of an error response from the Pterodactyl API.
type errorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

// APIError represents an error returned from the Pterodactyl API.
type APIError struct {
	// StatusCode is the HTTP status code from the response.
	StatusCode int
	// RateLimit is the rate limit information from the response headers.
	RateLimit transport.RateLimitInfo
	// Errors is a slice of specific error details from the API.
	Errors []ErrorDetail
	// Method is the HTTP method that was used.
	Method string
	// URL is the URL that was requested (with sensitive info redacted).
	URL string
	// raw is the raw response body, for debugging.
	raw []byte
}

// Error implements the error interface. It provides detailed context about the API error.
func (e *APIError) Error() string {
	var parts []string

	// Add method and URL
	if e.Method != "" && e.URL != "" {
		parts = append(parts, fmt.Sprintf("%s %s", e.Method, redactURL(e.URL)))
	}

	// Add status code
	parts = append(parts, fmt.Sprintf("status %d", e.StatusCode))

	// Add error details
	if len(e.Errors) > 0 {
		var errorMsgs []string
		for _, err := range e.Errors {
			msg := err.Detail
			if err.Source != nil && err.Source.Field != "" {
				msg = fmt.Sprintf("%s (field: %s)", msg, err.Source.Field)
			}
			errorMsgs = append(errorMsgs, msg)
		}
		parts = append(parts, strings.Join(errorMsgs, "; "))
	}

	return fmt.Sprintf("pterodactyl: %s", strings.Join(parts, ": "))
}

// redactURL removes sensitive information from URLs for error messages.
func redactURL(rawURL string) string {
	// Redact any tokens or sensitive query parameters
	// For now, just return as-is since our URLs don't contain sensitive data in path
	return rawURL
}

// CheckResponse checks the API response for errors, and returns one if present.
// A response is considered an error if it has a status code outside the 2xx range.
func CheckResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	apiErr := &APIError{
		StatusCode: r.StatusCode,
		RateLimit:  transport.ParseRateLimit(r),
		Method:     r.Request.Method,
		URL:        r.Request.URL.String(),
	}

	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		apiErr.raw = data
		errResponse := &errorResponse{}
		if err := json.Unmarshal(data, errResponse); err == nil {
			apiErr.Errors = errResponse.Errors
		}
	}

	return apiErr
}
