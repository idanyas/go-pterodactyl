package appapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/errors"
)

// mockRequester implements the requester.Requester interface for testing
type mockRequester struct {
	requests     []mockRequest
	responses    []mockResponse
	currentIndex int
}

type mockRequest struct {
	method   string
	endpoint string
	body     []byte
	options  *api.PaginationOptions
}

type mockResponse struct {
	statusCode int
	body       []byte
	err        error
}

func (m *mockRequester) NewRequest(ctx context.Context, method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error) {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	m.requests = append(m.requests, mockRequest{
		method:   method,
		endpoint: endpoint,
		body:     bodyBytes,
		options:  options,
	})

	// Create a minimal request for testing
	req, _ := http.NewRequestWithContext(ctx, method, "http://test.com"+endpoint, bytes.NewReader(bodyBytes))
	return req, nil
}

func (m *mockRequester) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	if m.currentIndex >= len(m.responses) {
		return nil, fmt.Errorf("no more mock responses available")
	}

	response := m.responses[m.currentIndex]
	m.currentIndex++

	if response.err != nil {
		return nil, response.err
	}

	// Create a mock response
	resp := &http.Response{
		StatusCode: response.statusCode,
		Body:       io.NopCloser(bytes.NewReader(response.body)),
	}

	// If we have a target to decode into and the response is successful
	if v != nil && response.statusCode >= 200 && response.statusCode < 300 {
		if err := json.NewDecoder(bytes.NewReader(response.body)).Decode(v); err != nil {
			return nil, err
		}
	}

	// If it's an error response, create an APIError
	if response.statusCode >= 400 {
		apiErr := &errors.APIError{HTTPStatusCode: response.statusCode}
		if len(response.body) > 0 {
			err := json.NewDecoder(bytes.NewReader(response.body)).Decode(apiErr)
			if err != nil {
				return nil, err
			}
		}
		return nil, apiErr
	}

	return resp, nil
}

func TestAllocationsService_List(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodeID         int
		options        *api.PaginationOptions
		mockResponse   mockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:   "Successful list with pagination",
			nodeID: 1,
			options: &api.PaginationOptions{
				Page:    1,
				PerPage: 10,
			},
			mockResponse: mockResponse{
				statusCode: 200,
				body: []byte(`{
					"object": "list",
					"data": [
						{
							"object": "allocation",
							"attributes": {
								"id": 1,
								"ip": "192.168.1.1",
								"port": 25565,
								"assigned": false
							}
						},
						{
							"object": "allocation",
							"attributes": {
								"id": 2,
								"ip": "192.168.1.1",
								"port": 25566,
								"assigned": true
							}
						}
					],
					"meta": {
						"pagination": {
							"total": 2,
							"count": 2,
							"per_page": 10,
							"current_page": 1,
							"total_pages": 1
						}
					}
				}`),
			},
			expectedError:  false,
			expectedCount:  2,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/1/allocations",
		},
		{
			name:    "Successful list without pagination",
			nodeID:  2,
			options: nil,
			mockResponse: mockResponse{
				statusCode: 200,
				body: []byte(`{
					"object": "list",
					"data": [],
					"meta": {
						"pagination": {
							"total": 0,
							"count": 0,
							"per_page": 100,
							"current_page": 1,
							"total_pages": 0
						}
					}
				}`),
			},
			expectedError:  false,
			expectedCount:  0,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/2/allocations",
		},
		{
			name:    "API error response",
			nodeID:  3,
			options: &api.PaginationOptions{Page: 1},
			mockResponse: mockResponse{
				statusCode: 404,
				body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested node could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/3/allocations",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{
				responses: []mockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			allocations, meta, err := service.List(context.Background(), tc.options)

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check response expectations
			if len(allocations) != tc.expectedCount {
				t.Errorf("expected %d allocations, got %d", tc.expectedCount, len(allocations))
			}

			// Check request expectations
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

			// If we have pagination options, verify they were passed correctly
			if tc.options != nil {
				if req.options == nil {
					t.Error("expected pagination options to be passed")
				} else if req.options.Page != tc.options.Page {
					t.Errorf("expected page %d, got %d", tc.options.Page, req.options.Page)
				}
			}

			// Verify meta is present for successful responses
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}
		})
	}
}

func TestAllocationsService_ListAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nodeID        int
		mockResponses []mockResponse
		expectedError bool
		expectedCount int
		expectedCalls int
	}{
		{
			name:   "Single page of results",
			nodeID: 1,
			mockResponses: []mockResponse{
				{
					statusCode: 200,
					body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "allocation",
								"attributes": {
									"id": 1,
									"ip": "192.168.1.1",
									"port": 25565,
									"assigned": false
								}
							}
						],
						"meta": {
							"pagination": {
								"total": 1,
								"count": 1,
								"per_page": 100,
								"current_page": 1,
								"total_pages": 1
							}
						}
					}`),
				},
			},
			expectedError: false,
			expectedCount: 1,
			expectedCalls: 1,
		},
		{
			name:   "Multiple pages of results",
			nodeID: 2,
			mockResponses: []mockResponse{
				{
					statusCode: 200,
					body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "allocation",
								"attributes": {
									"id": 1,
									"ip": "192.168.1.1",
									"port": 25565,
									"assigned": false
								}
							}
						],
						"meta": {
							"pagination": {
								"total": 2,
								"count": 1,
								"per_page": 1,
								"current_page": 1,
								"total_pages": 2
							}
						}
					}`),
				},
				{
					statusCode: 200,
					body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "allocation",
								"attributes": {
									"id": 2,
									"ip": "192.168.1.1",
									"port": 25566,
									"assigned": true
								}
							}
						],
						"meta": {
							"pagination": {
								"total": 2,
								"count": 1,
								"per_page": 1,
								"current_page": 2,
								"total_pages": 2
							}
						}
					}`),
				},
			},
			expectedError: false,
			expectedCount: 2,
			expectedCalls: 2,
		},
		{
			name:   "API error on first page",
			nodeID: 3,
			mockResponses: []mockResponse{
				{
					statusCode: 500,
					body: []byte(`{
						"errors": [{
							"code": "InternalServerError",
							"status": "500",
							"detail": "Internal server error"
						}]
					}`),
				},
			},
			expectedError: true,
			expectedCalls: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{
				responses: tc.mockResponses,
			}

			service := newAllocationsService(mock, tc.nodeID)

			allocations, err := service.ListAll(context.Background())

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check response expectations
			if len(allocations) != tc.expectedCount {
				t.Errorf("expected %d allocations, got %d", tc.expectedCount, len(allocations))
			}

			// Check number of API calls
			if len(mock.requests) != tc.expectedCalls {
				t.Errorf("expected %d API calls, got %d", tc.expectedCalls, len(mock.requests))
			}

			// Verify all requests were to the correct endpoint
			for i, req := range mock.requests {
				expectedPath := fmt.Sprintf("/api/application/nodes/%d/allocations", tc.nodeID)
				if req.endpoint != expectedPath {
					t.Errorf("request %d: expected path %s, got %s", i, expectedPath, req.endpoint)
				}

				if req.method != "GET" {
					t.Errorf("request %d: expected method GET, got %s", i, req.method)
				}
			}
		})
	}
}

func TestAllocationsService_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nodeID        int
		options       api.AllocationCreateOptions
		mockResponse  mockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:   "Successful creation",
			nodeID: 1,
			options: api.AllocationCreateOptions{
				IP:    "192.168.1.1",
				Ports: []string{"25565", "25566", "25567"},
			},
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(""),
			},
			expectedError: false,
			expectedBody:  `{"ip":"192.168.1.1","ports":["25565","25566","25567"]}`,
		},
		{
			name:   "Single port creation",
			nodeID: 2,
			options: api.AllocationCreateOptions{
				IP:    "10.0.0.1",
				Ports: []string{"8080"},
			},
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(""),
			},
			expectedError: false,
			expectedBody:  `{"ip":"10.0.0.1","ports":["8080"]}`,
		},
		{
			name:   "API error response",
			nodeID: 3,
			options: api.AllocationCreateOptions{
				IP:    "invalid-ip",
				Ports: []string{"99999"},
			},
			mockResponse: mockResponse{
				statusCode: 422,
				body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "The given data was invalid."
					}]
				}`),
			},
			expectedError: true,
			expectedBody:  `{"ip":"invalid-ip","ports":["99999"]}`,
		},
		{
			name:   "Network error",
			nodeID: 4,
			options: api.AllocationCreateOptions{
				IP:    "192.168.1.1",
				Ports: []string{"25565"},
			},
			mockResponse: mockResponse{
				err: fmt.Errorf("network timeout"),
			},
			expectedError: true,
			expectedBody:  `{"ip":"192.168.1.1","ports":["25565"]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{
				responses: []mockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			err := service.Create(context.Background(), tc.options)

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check request expectations
			if len(mock.requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.requests))
			}

			req := mock.requests[0]
			if req.method != "POST" {
				t.Errorf("expected method POST, got %s", req.method)
			}

			expectedPath := fmt.Sprintf("/api/application/nodes/%d/allocations", tc.nodeID)
			if req.endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.endpoint)
			}

			// Check request body
			bodyStr := strings.TrimSpace(string(req.body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected body %s, got %s", tc.expectedBody, bodyStr)
			}
		})
	}
}

func TestAllocationsService_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodeID         int
		allocationID   int
		mockResponse   mockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:         "Successful deletion",
			nodeID:       1,
			allocationID: 123,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(""),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/1/allocations/123",
		},
		{
			name:         "Allocation not found",
			nodeID:       2,
			allocationID: 999,
			mockResponse: mockResponse{
				statusCode: 404,
				body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested allocation could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/2/allocations/999",
		},
		{
			name:         "Allocation in use",
			nodeID:       3,
			allocationID: 456,
			mockResponse: mockResponse{
				statusCode: 422,
				body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "Cannot delete allocation that is currently assigned to a server."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/3/allocations/456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{
				responses: []mockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			err := service.Delete(context.Background(), tc.allocationID)

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check request expectations
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

func TestAllocationsService_Integration(t *testing.T) {
	t.Parallel()

	// Test that the service correctly handles the nodeID in all operations
	nodeID := 42
	mock := &mockRequester{
		responses: []mockResponse{
			// List response
			{
				statusCode: 200,
				body: []byte(`{
					"object": "list",
					"data": [],
					"meta": {
						"pagination": {
							"total": 0,
							"count": 0,
							"per_page": 100,
							"current_page": 1,
							"total_pages": 0
						}
					}
				}`),
			},
			// ListAll response
			{
				statusCode: 200,
				body: []byte(`{
					"object": "list",
					"data": [],
					"meta": {
						"pagination": {
							"total": 0,
							"count": 0,
							"per_page": 100,
							"current_page": 1,
							"total_pages": 0
						}
					}
				}`),
			},
			// Create response
			{
				statusCode: 204,
				body:       []byte(""),
			},
			// Delete response
			{
				statusCode: 204,
				body:       []byte(""),
			},
		},
	}

	service := newAllocationsService(mock, nodeID)

	// Test List
	_, _, err := service.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	// Test ListAll
	_, err = service.ListAll(context.Background())
	if err != nil {
		t.Fatalf("ListAll failed: %v", err)
	}

	// Test Create
	createOptions := api.AllocationCreateOptions{
		IP:    "192.168.1.1",
		Ports: []string{"25565"},
	}
	err = service.Create(context.Background(), createOptions)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test Delete
	err = service.Delete(context.Background(), 123)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify all requests used the correct nodeID
	expectedBasePath := fmt.Sprintf("/api/application/nodes/%d/allocations", nodeID)
	for i, req := range mock.requests {
		if !strings.HasPrefix(req.endpoint, expectedBasePath) {
			t.Errorf("request %d: expected endpoint to start with %s, got %s", i, expectedBasePath, req.endpoint)
		}
	}

	// Verify we made exactly 4 requests
	if len(mock.requests) != 4 {
		t.Errorf("expected 4 requests, got %d", len(mock.requests))
	}
}

func TestAllocationsService_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nodeID        int
		allocationID  int
		options       api.AllocationCreateOptions
		mockResponse  mockResponse
		expectedError bool
		description   string
	}{
		{
			name:         "Zero node ID",
			nodeID:       0,
			allocationID: 123,
			mockResponse: mockResponse{
				statusCode: 204,
				body:       []byte(""),
			},
			expectedError: false,
			description:   "Should handle zero node ID gracefully",
		},
		{
			name:         "Negative allocation ID",
			nodeID:       1,
			allocationID: -1,
			mockResponse: mockResponse{
				statusCode: 400,
				body: []byte(`{
					"errors": [{
						"code": "BadRequestHttpException",
						"status": "400",
						"detail": "Invalid allocation ID"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle negative allocation ID with error",
		},
		{
			name:   "Empty ports array",
			nodeID: 1,
			options: api.AllocationCreateOptions{
				IP:    "192.168.1.1",
				Ports: []string{},
			},
			mockResponse: mockResponse{
				statusCode: 422,
				body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "At least one port must be specified"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle empty ports array with validation error",
		},
		{
			name:   "Empty IP address",
			nodeID: 1,
			options: api.AllocationCreateOptions{
				IP:    "",
				Ports: []string{"25565"},
			},
			mockResponse: mockResponse{
				statusCode: 422,
				body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "IP address is required"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle empty IP address with validation error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRequester{
				responses: []mockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			var err error
			if tc.options.IP != "" || len(tc.options.Ports) > 0 {
				// Test Create
				err = service.Create(context.Background(), tc.options)
			} else {
				// Test Delete
				err = service.Delete(context.Background(), tc.allocationID)
			}

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error for %s but got none", tc.description)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s: %v", tc.description, err)
				}
			}
		})
	}
}
