package appapi

import (
	"context"
	"testing"

	"github.com/davidarkless/go-pterodactyl/api"
)

func TestServersService_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		options        api.PaginationOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:    "Successful list with pagination",
			options: api.PaginationOptions{Page: 1, PerPage: 10},
			mockResponse: mockResponse{
				statusCode: 200,
				body: []byte(`{"object": "list", "data": [
					{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "Server1", "description": "desc1", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}
				], "meta": {"pagination": {"total": 1, "count": 1, "per_page": 10, "current_page": 1, "total_pages": 1}}}`),
			},
			expectedError:  false,
			expectedCount:  1,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
		{
			name:    "API error response",
			options: api.PaginationOptions{Page: 1},
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			servers, meta, err := service.List(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(servers) != tc.expectedCount {
				t.Errorf("expected %d servers, got %d", tc.expectedCount, len(servers))
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
			if tc.expectedCount > 0 && len(servers) > 0 {
				server := servers[0]
				if server.ID == 0 {
					t.Error("expected server ID to be non-zero")
				}
				if server.Name == "" {
					t.Error("expected server name to be non-empty")
				}
			}
		})
	}
}

func TestServersService_ListAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		mockResponse   mockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful list all",
			mockResponse: mockResponse{
				statusCode: 200,
				body: []byte(`{"object": "list", "data": [
					{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "Server1", "description": "desc1", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}},
					{"object": "server", "attributes": {"id": 2, "uuid": "uuid-2", "identifier": "id2", "name": "Server2", "description": "desc2", "suspended": true, "user": 2, "node": 1, "allocation": 2, "nest": 1, "egg": 1, "created_at": "2023-01-02T00:00:00Z"}}
				], "meta": {"pagination": {"total": 2, "count": 2, "per_page": 100, "current_page": 1, "total_pages": 1}}}`),
			},
			expectedError:  false,
			expectedCount:  2,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
		{
			name: "Empty list",
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "list", "data": [], "meta": {"pagination": {"total": 0, "count": 0, "per_page": 100, "current_page": 1, "total_pages": 0}}}`),
			},
			expectedError:  false,
			expectedCount:  0,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
		{
			name: "API error response",
			mockResponse: mockResponse{
				statusCode: 500,
				body:       []byte(`{"errors": [{"code": "InternalServerError", "status": "500", "detail": "Internal server error."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			servers, err := service.ListAll(context.Background())
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(servers) != tc.expectedCount {
				t.Errorf("expected %d servers, got %d", tc.expectedCount, len(servers))
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
			if tc.expectedCount > 0 && len(servers) > 0 {
				server := servers[0]
				if server.ID == 0 {
					t.Error("expected server ID to be non-zero")
				}
				if server.Name == "" {
					t.Error("expected server name to be non-empty")
				}
			}
		})
	}
}

func TestServersService_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful get",
			serverID: 1,
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/1",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.Get(context.Background(), tc.serverID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
			if server.ID != tc.serverID {
				t.Errorf("expected server ID %d, got %d", tc.serverID, server.ID)
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
		})
	}
}

func TestServersService_GetExternal(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		externalID     string
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:       "Successful get external",
			externalID: "external-123",
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/external/external-123",
		},
		{
			name:       "External server not found",
			externalID: "nonexistent",
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "External server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/external/nonexistent",
		},
		{
			name:       "URL encoding test",
			externalID: "external/with/slashes",
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/external/external%2Fwith%2Fslashes",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.GetExternal(context.Background(), tc.externalID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
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
		})
	}
}

func TestServersService_Create(t *testing.T) {
	t.Parallel()
	description := "Test server description"
	nodeID := 1
	allocation := &struct {
		Default    int   `json:"default"`
		Additional []int `json:"additional,omitempty"`
	}{
		Default: 1,
	}

	testCases := []struct {
		name           string
		options        api.ServerCreateOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful create",
			options: api.ServerCreateOptions{
				Name:        "TestServer",
				User:        1,
				NodeID:      &nodeID,
				Allocation:  allocation,
				Nest:        1,
				Egg:         1,
				Description: &description,
				Limits: api.ServerLimits{
					Memory: 1024,
					Swap:   512,
					Disk:   10000,
					IO:     500,
					CPU:    100,
				},
				FeatureLimits: api.ServerFeatureLimits{
					Databases:   5,
					Allocations: 1,
					Backups:     2,
				},
			},
			mockResponse: mockResponse{
				statusCode: 201,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test server description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers",
		},
		{
			name: "API error response",
			options: api.ServerCreateOptions{
				Name:       "TestServer",
				User:       1,
				NodeID:     &nodeID,
				Allocation: allocation,
				Nest:       1,
				Egg:        1,
				Limits: api.ServerLimits{
					Memory: 1024,
					Swap:   512,
					Disk:   10000,
					IO:     500,
					CPU:    100,
				},
				FeatureLimits: api.ServerFeatureLimits{
					Databases:   5,
					Allocations: 1,
					Backups:     2,
				},
			},
			mockResponse: mockResponse{
				statusCode: 422,
				body:       []byte(`{"errors": [{"code": "ValidationException", "status": "422", "detail": "Validation failed."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.Create(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
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
			// Verify request body contains expected data
			if req.body == nil {
				t.Error("expected request body to be non-nil")
			}
		})
	}
}

func TestServersService_UpdateDetails(t *testing.T) {
	t.Parallel()
	description := "Updated description"
	externalID := "external-123"

	testCases := []struct {
		name           string
		serverID       int
		options        api.ServerUpdateDetailsOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful update details",
			serverID: 1,
			options: api.ServerUpdateDetailsOptions{
				Name:        "UpdatedServer",
				Description: &description,
				ExternalID:  &externalID,
			},
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "UpdatedServer", "description": "Updated description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/details",
		},
		{
			name:     "Server not found",
			serverID: 999,
			options: api.ServerUpdateDetailsOptions{
				Name: "UpdatedServer",
			},
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/999/details",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.UpdateDetails(context.Background(), tc.serverID, tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
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
		})
	}
}

func TestServersService_UpdateBuild(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		options        api.ServerUpdateBuildOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful update build",
			serverID: 1,
			options: api.ServerUpdateBuildOptions{
				Allocation: 2,
				Memory:     1024,
				Swap:       512,
				Disk:       10000,
				IO:         500,
				CPU:        100,
				FeatureLimits: api.ServerFeatureLimits{
					Databases:   5,
					Allocations: 1,
					Backups:     2,
				},
			},
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 2, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/build",
		},
		{
			name:     "Invalid build configuration",
			serverID: 1,
			options: api.ServerUpdateBuildOptions{
				Allocation: 1,
				Memory:     -1, // Invalid memory value
			},
			mockResponse: mockResponse{
				statusCode: 422,
				body:       []byte(`{"errors": [{"code": "ValidationException", "status": "422", "detail": "Invalid memory allocation."}]}`),
			},
			expectedError:  true,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/build",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.UpdateBuild(context.Background(), tc.serverID, tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
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
		})
	}
}

func TestServersService_UpdateStartup(t *testing.T) {
	t.Parallel()
	environment := map[string]string{
		"JAVA_MEMORY": "1024M",
		"PORT":        "25565",
	}

	testCases := []struct {
		name           string
		serverID       int
		options        api.ServerUpdateStartupOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful update startup",
			serverID: 1,
			options: api.ServerUpdateStartupOptions{
				Startup:     "java -jar server.jar",
				Environment: &environment,
				Egg:         1,
				Image:       "openjdk:11",
				SkipScripts: false,
			},
			mockResponse: mockResponse{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/startup",
		},
		{
			name:     "Invalid startup configuration",
			serverID: 1,
			options: api.ServerUpdateStartupOptions{
				Startup: "", // Empty startup command
				Egg:     1,
				Image:   "openjdk:11",
			},
			mockResponse: mockResponse{
				statusCode: 422,
				body:       []byte(`{"errors": [{"code": "ValidationException", "status": "422", "detail": "Startup command cannot be empty."}]}`),
			},
			expectedError:  true,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/startup",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.UpdateStartup(context.Background(), tc.serverID, tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
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
		})
	}
}

func TestServersService_Suspend(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful suspend",
			serverID: 1,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/suspend",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/suspend",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Suspend(context.Background(), tc.serverID)
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
		})
	}
}

func TestServersService_Unsuspend(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful unsuspend",
			serverID: 1,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/unsuspend",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/unsuspend",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Unsuspend(context.Background(), tc.serverID)
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
		})
	}
}

func TestServersService_Reinstall(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful reinstall",
			serverID: 1,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/reinstall",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/reinstall",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Reinstall(context.Background(), tc.serverID)
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
		})
	}
}

func TestServersService_Delete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		force          bool
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful delete",
			serverID: 1,
			force:    false,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/1",
		},
		{
			name:     "Successful force delete",
			serverID: 1,
			force:    true,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/1",
		},
		{
			name:     "Server not found",
			serverID: 999,
			force:    false,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{responses: []mockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Delete(context.Background(), tc.serverID, tc.force)
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
			// Check if force delete includes request body
			if tc.force && req.body == nil {
				t.Error("expected request body for force delete")
			}
		})
	}
}

func TestServersService_Databases(t *testing.T) {
	t.Parallel()
	service := NewServersService(&mockRequester{})

	// Test that Databases returns a non-nil DatabaseService
	databaseService := service.Databases(context.Background(), 1)
	if databaseService == nil {
		t.Fatal("expected DatabaseService to be non-nil")
	}

}

// Test data validation and edge cases
func TestServersService_DataValidation(t *testing.T) {
	t.Parallel()

	t.Run("Empty server list response", func(t *testing.T) {
		mock := &mockRequester{responses: []mockResponse{
			{
				statusCode: 200,
				body:       []byte(`{"object": "list", "data": [], "meta": {"pagination": {"total": 0, "count": 0, "per_page": 100, "current_page": 1, "total_pages": 0}}}`),
			},
		}}
		service := NewServersService(mock)
		servers, err := service.ListAll(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(servers) != 0 {
			t.Errorf("expected 0 servers, got %d", len(servers))
		}
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		mock := &mockRequester{responses: []mockResponse{
			{
				statusCode: 200,
				body:       []byte(`invalid json`),
			},
		}}
		service := NewServersService(mock)
		_, err := service.ListAll(context.Background())
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("Server with minimal data", func(t *testing.T) {
		mock := &mockRequester{responses: []mockResponse{
			{
				statusCode: 200,
				body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "MinimalServer", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
		}}
		service := NewServersService(mock)
		server, err := service.Get(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if server == nil {
			t.Fatal("expected server to be non-nil")
		}
		if server.ID != 1 {
			t.Errorf("expected server ID 1, got %d", server.ID)
		}
		if server.Name != "MinimalServer" {
			t.Errorf("expected server name 'MinimalServer', got '%s'", server.Name)
		}
	})
}

// Test constructor
func TestNewServersService(t *testing.T) {
	t.Parallel()
	mock := &mockRequester{}
	service := NewServersService(mock)
	if service == nil {
		t.Fatal("expected service to be non-nil")
	}
}
