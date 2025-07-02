package appapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/davidarkless/go-pterodactyl/api"
)

func TestUsersService_List(t *testing.T) {
	t.Parallel()

	mockUserListResponse := `{
		"object": "list",
		"data": [
			{
				"object": "user",
				"attributes": { "id": 1, "username": "testuser1", "email": "test1@example.com" }
			},
			{
				"object": "user",
				"attributes": { "id": 2, "username": "testuser2", "email": "test2@example.com" }
			}
		],
		"meta": {
			"pagination": { "total": 2, "count": 2, "per_page": 10, "current_page": 1, "total_pages": 1 }
		}
	}`

	testCases := []struct {
		name          string
		options       *api.PaginationOptions
		mockResponse  mockResponse
		expectedError bool
		expectedCount int
	}{
		{
			name:    "Successful list",
			options: &api.PaginationOptions{Page: 1},
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(mockUserListResponse),
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:    "API error",
			options: &api.PaginationOptions{Page: 1},
			mockResponse: mockResponse{
				statusCode: 500,
				body:       []byte(`{"errors": [{"code": "ServerException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			users, meta, err := service.List(context.Background(), tc.options)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(users) != tc.expectedCount {
				t.Errorf("expected %d users, got %d", tc.expectedCount, len(users))
			}
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}
			if mock.requests[0].endpoint != "/api/application/users" {
				t.Errorf("expected endpoint '/api/application/users', got '%s'", mock.requests[0].endpoint)
			}
		})
	}
}

func TestUsersService_Get(t *testing.T) {
	t.Parallel()

	mockUserResponse := `{
		"object": "user",
		"attributes": { "id": 123, "username": "getuser", "email": "get@example.com" }
	}`

	testCases := []struct {
		name          string
		userID        int
		mockResponse  mockResponse
		expectedError bool
		expectedID    int
	}{
		{
			name:   "Successful get",
			userID: 123,
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(mockUserResponse),
			},
			expectedError: false,
			expectedID:    123,
		},
		{
			name:   "User not found",
			userID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.Get(context.Background(), tc.userID)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.ID != tc.expectedID {
				t.Errorf("expected user ID %d, got %d", tc.expectedID, user.ID)
			}

			expectedPath := fmt.Sprintf("/api/application/users/%d", tc.userID)
			if mock.requests[0].endpoint != expectedPath {
				t.Errorf("expected endpoint '%s', got '%s'", expectedPath, mock.requests[0].endpoint)
			}
		})
	}
}

func TestUsersService_GetExternalID(t *testing.T) {
	t.Parallel()

	mockUserResponse := `{
		"object": "user",
		"attributes": { "id": 456, "external_id": "ext-123", "username": "extuser" }
	}`

	testCases := []struct {
		name           string
		externalID     string
		mockResponse   mockResponse
		expectedError  bool
		expectedUserID int
	}{
		{
			name:       "Successful get by external ID",
			externalID: "ext-123",
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(mockUserResponse),
			},
			expectedError:  false,
			expectedUserID: 456,
		},
		{
			name:       "User not found by external ID",
			externalID: "ext-999",
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.GetExternalID(context.Background(), tc.externalID)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.ID != tc.expectedUserID {
				t.Errorf("expected user ID %d, got %d", tc.expectedUserID, user.ID)
			}

			expectedPath := fmt.Sprintf("/api/application/users/external/%s", tc.externalID)
			if mock.requests[0].endpoint != expectedPath {
				t.Errorf("expected endpoint '%s', got '%s'", expectedPath, mock.requests[0].endpoint)
			}
		})
	}
}

func TestUsersService_Create(t *testing.T) {
	t.Parallel()

	options := api.UserCreateOptions{
		Username:  "newuser",
		Email:     "new@example.com",
		FirstName: "New",
		LastName:  "User",
		Password:  "supersecret",
	}

	mockCreatedUserResponse := `{
		"object": "user",
		"attributes": { "id": 1, "username": "newuser", "email": "new@example.com" }
	}`

	expectedBody, _ := json.Marshal(options)

	testCases := []struct {
		name          string
		options       api.UserCreateOptions
		mockResponse  mockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful creation",
			options: options,
			mockResponse: mockResponse{
				statusCode: 201, // 201 Created is typical
				body:       []byte(mockCreatedUserResponse),
			},
			expectedError: false,
			expectedBody:  string(expectedBody),
		},
		{
			name:    "Validation error",
			options: options,
			mockResponse: mockResponse{
				statusCode: 422,
				body:       []byte(`{"errors": [{"code": "ValidationException"}]}`),
			},
			expectedError: true,
			expectedBody:  string(expectedBody),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.Create(context.Background(), tc.options)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user == nil {
				t.Fatal("expected a user object but got nil")
			}
			if user.Username != tc.options.Username {
				t.Errorf("expected username '%s', got '%s'", tc.options.Username, user.Username)
			}

			req := mock.requests[0]
			if req.method != "POST" {
				t.Errorf("expected method POST, got %s", req.method)
			}
			if req.endpoint != "/api/application/users" {
				t.Errorf("expected endpoint '/api/application/users', got '%s'", req.endpoint)
			}
			if strings.TrimSpace(string(req.body)) != tc.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tc.expectedBody, string(req.body))
			}
		})
	}
}

func TestUsersService_Update(t *testing.T) {
	t.Parallel()

	userID := 123
	options := api.UserUpdateOptions{
		Username:  "updateduser",
		FirstName: "Updated",
	}

	mockUpdatedUserResponse := `{
		"object": "user",
		"attributes": { "id": 123, "username": "updateduser", "first_name": "Updated" }
	}`

	expectedBody, _ := json.Marshal(options)

	testCases := []struct {
		name          string
		userID        int
		options       api.UserUpdateOptions
		mockResponse  mockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful update",
			userID:  userID,
			options: options,
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(mockUpdatedUserResponse),
			},
			expectedError: false,
			expectedBody:  string(expectedBody),
		},
		{
			name:    "User not found",
			userID:  999,
			options: options,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
			expectedBody:  string(expectedBody),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.Update(context.Background(), tc.userID, tc.options)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.Username != tc.options.Username {
				t.Errorf("expected username '%s', got '%s'", tc.options.Username, user.Username)
			}

			req := mock.requests[0]
			if req.method != "PATCH" {
				t.Errorf("expected method PATCH, got %s", req.method)
			}
			expectedPath := fmt.Sprintf("/api/application/users/%d", tc.userID)
			if req.endpoint != expectedPath {
				t.Errorf("expected endpoint '%s', got '%s'", expectedPath, req.endpoint)
			}
			if strings.TrimSpace(string(req.body)) != tc.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tc.expectedBody, string(req.body))
			}
		})
	}
}

func TestUsersService_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        int
		mockResponse  mockResponse
		expectedError bool
	}{
		{
			name:   "Successful deletion",
			userID: 123,
			mockResponse: mockResponse{
				statusCode: 204, // No Content
				body:       []byte(""),
			},
			expectedError: false,
		},
		{
			name:   "User not found",
			userID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
		},
		{
			name:   "Cannot delete user with servers",
			userID: 456,
			mockResponse: mockResponse{
				statusCode: 400, // Bad Request
				body:       []byte(`{"errors": [{"code": "HasActiveServersException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			err := service.Delete(context.Background(), tc.userID)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			req := mock.requests[0]
			if req.method != "DELETE" {
				t.Errorf("expected method DELETE, got %s", req.method)
			}
			expectedPath := fmt.Sprintf("/api/application/users/%d", tc.userID)
			if req.endpoint != expectedPath {
				t.Errorf("expected endpoint '%s', got '%s'", expectedPath, req.endpoint)
			}
		})
	}
}
