package appapi

import (
	"context"
	"testing"
	"time"

	"github.com/davidarkless/go-pterodactyl/api"
)

func TestNestsService_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		options        *api.PaginationOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:    "Successful list with pagination",
			options: &api.PaginationOptions{Page: 1, PerPage: 10},
			mockResponse: mockResponse{
				statusCode: 200,
				body: []byte(`{
					"object": "list",
					"data": [
						{"object": "nest", "attributes": {"id": 1, "uuid": "uuid-1", "author": "author1", "name": "Nest1", "description": "desc1", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}},
						{"object": "nest", "attributes": {"id": 2, "uuid": "uuid-2", "author": "author2", "name": "Nest2", "description": "desc2", "created_at": "2023-01-02T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z"}}
					],
					"meta": {"pagination": {"total": 2, "count": 2, "per_page": 10, "current_page": 1, "total_pages": 1}}
				}`),
			},
			expectedError:  false,
			expectedCount:  2,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests",
		},
		{
			name:    "Successful list without pagination",
			options: nil,
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "list", "data": [], "meta": {"pagination": {"total": 0, "count": 0, "per_page": 100, "current_page": 1, "total_pages": 0}}}`),
			},
			expectedError:  false,
			expectedCount:  0,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests",
		},
		{
			name:    "API error response",
			options: &api.PaginationOptions{Page: 1},
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewNestsService(mock)
			nests, meta, err := service.List(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(nests) != tc.expectedCount {
				t.Errorf("expected %d nests, got %d", tc.expectedCount, len(nests))
			}
			if len(mock.requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.requests))
			}
			req := mock.requests[0]
			if req.method != tc.expectedMethod {
				t.Errorf("expected method %s, got %s", tc.expectedMethod, req.method)
			}
			if req.endpoint != tc.expectedPath {
				t.Errorf("expected path %s, got %s", tc.expectedPath, req.endpoint)
			}
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}
			if tc.expectedCount > 0 && len(nests) > 0 {
				nest := nests[0]
				if nest.ID == 0 {
					t.Error("expected nest ID to be non-zero")
				}
				if nest.Name == "" {
					t.Error("expected nest name to be non-empty")
				}
			}
		})
	}
}

func TestNestsService_ListAll(t *testing.T) {
	t.Parallel()
	mock := &mockRequester{responses: []mockResponse{{statusCode: 200, body: []byte(`{"object": "list", "data": [{"object": "nest", "attributes": {"id": 1, "uuid": "uuid-1", "author": "author1", "name": "Nest1", "description": "desc1", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}], "meta": {"pagination": {"total": 1, "count": 1, "per_page": 100, "current_page": 1, "total_pages": 1}}}}`)}}}
	service := NewNestsService(mock)
	nests, err := service.ListAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nests) != 1 {
		t.Errorf("expected 1 nest, got %d", len(nests))
	}
	if nests[0].Name != "Nest1" {
		t.Errorf("expected name 'Nest1', got '%s'", nests[0].Name)
	}
}

func TestNestsService_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		id             int
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful get",
			id:   1,
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "nest", "attributes": {"id": 1, "uuid": "uuid-1", "author": "author1", "name": "Nest1", "description": "desc1", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests/1",
		},
		{
			name: "Nest not found",
			id:   999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewNestsService(mock)
			nest, err := service.Get(context.Background(), tc.id)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(mock.requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.requests))
			}
			req := mock.requests[0]
			if req.method != tc.expectedMethod {
				t.Errorf("expected method %s, got %s", tc.expectedMethod, req.method)
			}
			if req.endpoint != tc.expectedPath {
				t.Errorf("expected path %s, got %s", tc.expectedPath, req.endpoint)
			}
			if nest == nil {
				t.Error("expected nest to be non-nil")
			} else {
				if nest.ID != tc.id {
					t.Errorf("expected nest ID %d, got %d", tc.id, nest.ID)
				}
				if nest.Name == "" {
					t.Error("expected nest name to be non-empty")
				}
			}
		})
	}
}

func TestNestsService_Eggs(t *testing.T) {
	t.Parallel()
	mock := &mockRequester{}
	service := NewNestsService(mock)
	eggsService := service.Eggs(42)
	if eggsService == nil {
		t.Fatal("expected eggsService to be non-nil")
	}
}

func TestNestsService_DataValidation(t *testing.T) {
	t.Parallel()
	mock := &mockRequester{responses: []mockResponse{{statusCode: 200, body: []byte(`{"object": "nest", "attributes": {"id": 1, "uuid": "uuid-1", "author": "author1", "name": "Nest1", "description": "desc1", "created_at": "2023-01-01T12:00:00Z", "updated_at": "2023-01-02T12:00:00Z"}}`)}}}
	service := NewNestsService(mock)
	nest, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nest.ID != 1 {
		t.Errorf("expected ID 1, got %d", nest.ID)
	}
	if nest.UUID != "uuid-1" {
		t.Errorf("expected UUID 'uuid-1', got '%s'", nest.UUID)
	}
	if nest.Author != "author1" {
		t.Errorf("expected Author 'author1', got '%s'", nest.Author)
	}
	if nest.Name != "Nest1" {
		t.Errorf("expected Name 'Nest1', got '%s'", nest.Name)
	}
	if nest.Description != "desc1" {
		t.Errorf("expected Description 'desc1', got '%s'", nest.Description)
	}
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	if !nest.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", expectedCreatedAt, nest.CreatedAt)
	}
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2023-01-02T12:00:00Z")
	if !nest.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", expectedUpdatedAt, nest.UpdatedAt)
	}
}
