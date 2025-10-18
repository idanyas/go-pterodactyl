package validation

import (
	"strings"
	"testing"
)

type testStruct struct {
	Name  string `validate:"required,min=3"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=120"`
	URL   string `validate:"omitempty,url"`
}

func TestValidate_ValidStruct(t *testing.T) {
	valid := testStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		URL:   "https://example.com",
	}

	err := Validate(valid)
	if err != nil {
		t.Errorf("Validate() failed for valid struct: %v", err)
	}
}

func TestValidate_InvalidStruct(t *testing.T) {
	tests := []struct {
		name        string
		input       testStruct
		expectError bool
		errorField  string
	}{
		{
			name: "missing required name",
			input: testStruct{
				Email: "john@example.com",
				Age:   30,
			},
			expectError: true,
			errorField:  "Name",
		},
		{
			name: "name too short",
			input: testStruct{
				Name:  "Jo",
				Email: "john@example.com",
				Age:   30,
			},
			expectError: true,
			errorField:  "Name",
		},
		{
			name: "invalid email",
			input: testStruct{
				Name:  "John Doe",
				Email: "not-an-email",
				Age:   30,
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "age too low",
			input: testStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   -1,
			},
			expectError: true,
			errorField:  "Age",
		},
		{
			name: "age too high",
			input: testStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   121,
			},
			expectError: true,
			errorField:  "Age",
		},
		{
			name: "invalid URL",
			input: testStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
				URL:   "not a url",
			},
			expectError: true,
			errorField:  "URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("expected validation error, got nil")
					return
				}

				valErr, ok := err.(*ValidationError)
				if !ok {
					t.Errorf("expected *ValidationError, got %T", err)
					return
				}

				// Check that the expected field is in the errors
				found := false
				for _, fieldErr := range valErr.Errors {
					if fieldErr.Field == tt.errorField {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("expected error for field '%s', but it was not found in errors: %v", tt.errorField, valErr.Errors)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestValidate_NilInput(t *testing.T) {
	err := Validate(nil)
	if err != nil {
		t.Errorf("Validate(nil) should return nil, got: %v", err)
	}
}

func TestValidationError_Error(t *testing.T) {
	// Test with actual validation error
	invalid := testStruct{
		Name:  "Jo", // too short
		Email: "invalid-email",
		Age:   -5,
	}

	err := Validate(invalid)
	if err == nil {
		t.Fatal("expected validation error")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "validation failed") {
		t.Errorf("error message should contain 'validation failed', got: %s", errMsg)
	}

	// Should mention at least one field
	if !strings.Contains(errMsg, "Name") && !strings.Contains(errMsg, "Email") && !strings.Contains(errMsg, "Age") {
		t.Errorf("error message should mention failing fields, got: %s", errMsg)
	}
}

func TestValidate_OmitEmpty(t *testing.T) {
	// URL is optional (omitempty), so this should pass
	valid := testStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		URL:   "", // empty is ok with omitempty
	}

	err := Validate(valid)
	if err != nil {
		t.Errorf("Validate() failed for struct with empty optional field: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	// Multiple validation errors
	invalid := testStruct{
		Name:  "Jo",           // too short
		Email: "not-an-email", // invalid email
		Age:   200,            // too high
		URL:   "invalid-url",  // invalid URL
	}

	err := Validate(invalid)
	if err == nil {
		t.Fatal("expected validation error")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}

	if len(valErr.Errors) < 2 {
		t.Errorf("expected at least 2 validation errors, got %d", len(valErr.Errors))
	}
}
