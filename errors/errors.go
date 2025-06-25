package errors

import (
	"fmt"
	"strings"
)

type PterodactylError struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type APIError struct {
	HTTPStatusCode int
	Errors         []PterodactylError
}

// Error implements the error interface, providing a user-friendly error message.
func (e *APIError) Error() string {
	var errorDetails []string
	for _, pErr := range e.Errors {
		errorDetails = append(errorDetails, pErr.Detail)
	}
	return fmt.Sprintf("pterodactyl: API error (status %d): %s",
		e.HTTPStatusCode, strings.Join(errorDetails, ", "))
}
