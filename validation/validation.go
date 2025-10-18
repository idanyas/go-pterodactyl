// Package validation provides input validation utilities for API requests.
package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// validate is the singleton validator instance
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

// Validate validates a struct using struct tags.
// Returns a validation error if the struct does not pass validation,
// or nil if validation succeeds.
//
// Example:
//
//	type CreateUserRequest struct {
//	    Email    string `validate:"required,email"`
//	    Username string `validate:"required,min=3,max=20"`
//	    Age      int    `validate:"gte=0,lte=120"`
//	}
//
//	req := CreateUserRequest{Email: "invalid", Username: "ab"}
//	if err := validation.Validate(req); err != nil {
//	    // Handle validation error
//	}
func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	// Convert validator errors to a more readable format
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return NewValidationError(validationErrors)
}

// ValidationError represents one or more validation failures.
type ValidationError struct {
	Errors []FieldError
}

// FieldError represents a validation failure for a specific field.
type FieldError struct {
	Field   string
	Tag     string
	Value   interface{}
	Message string
}

// NewValidationError creates a ValidationError from validator.ValidationErrors.
func NewValidationError(errs validator.ValidationErrors) *ValidationError {
	fieldErrors := make([]FieldError, len(errs))
	for i, err := range errs {
		fieldErrors[i] = FieldError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   err.Value(),
			Message: formatErrorMessage(err),
		}
	}
	return &ValidationError{Errors: fieldErrors}
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation error"
	}

	var sb strings.Builder
	sb.WriteString("validation failed:")
	for _, err := range e.Errors {
		sb.WriteString("\n  - ")
		sb.WriteString(err.Message)
	}
	return sb.String()
}

// formatErrorMessage creates a human-readable error message for a validation error.
func formatErrorMessage(err validator.FieldError) string {
	field := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", field)
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email address", field)
	case "min":
		return fmt.Sprintf("field '%s' must be at least %s characters long", field, err.Param())
	case "max":
		return fmt.Sprintf("field '%s' must be at most %s characters long", field, err.Param())
	case "gte":
		return fmt.Sprintf("field '%s' must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("field '%s' must be less than or equal to %s", field, err.Param())
	case "gt":
		return fmt.Sprintf("field '%s' must be greater than %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("field '%s' must be less than %s", field, err.Param())
	case "len":
		return fmt.Sprintf("field '%s' must be exactly %s characters long", field, err.Param())
	case "oneof":
		return fmt.Sprintf("field '%s' must be one of [%s]", field, err.Param())
	case "url":
		return fmt.Sprintf("field '%s' must be a valid URL", field)
	case "ip":
		return fmt.Sprintf("field '%s' must be a valid IP address", field)
	case "ipv4":
		return fmt.Sprintf("field '%s' must be a valid IPv4 address", field)
	case "ipv6":
		return fmt.Sprintf("field '%s' must be a valid IPv6 address", field)
	default:
		return fmt.Sprintf("field '%s' failed validation '%s'", field, err.Tag())
	}
}
