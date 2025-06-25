package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type allocationsService struct {
	client requester.Requester
	nodeID int
}

func newAllocationsService(client requester.Requester, nodeID int) *allocationsService {
	return &allocationsService{client: client, nodeID: nodeID}
}

func (s *allocationsService) List(options *api.PaginationOptions) ([]*api.Allocation, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations", s.nodeID)
	req, err := s.client.NewRequest("GET", endpoint, nil, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create allocation list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Allocation]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Allocation, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *allocationsService) ListAll() ([]*api.Allocation, error) {
	// Start with the first page and a large-ish slice capacity.
	allAllocations := make([]*api.Allocation, 0, 100)
	options := &api.PaginationOptions{PerPage: 100, Page: 1} // Use max per_page for efficiency

	for {
		// Fetch a page of allocations.
		allocations, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}

		// Add the fetched allocations to our master slice.
		allAllocations = append(allAllocations, allocations...)

		// Check if we've reached the last page.
		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break // Exit the loop
		}

		options.Page++
	}

	return allAllocations, nil
}

func (s *allocationsService) Create(options api.AllocationCreateOptions) error {
	// Marshal the options struct into JSON.
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal create allocation options: %w", err)
	}

	// Construct the endpoint for creating allocations on this node.
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations", s.nodeID)
	req, err := s.client.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create new allocation request: %w", err)
	}

	// Execute the request. We pass `nil` for the decoding target because we
	// expect a 204 No Content response.
	_, err = s.client.Do(req, nil)
	return err
}

// Delete deletes a specific allocation from the configured node.
func (s *allocationsService) Delete(allocationID int) error {
	// Construct the specific endpoint for the allocation to be deleted.
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations/%d", s.nodeID, allocationID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete allocation request: %w", err)
	}

	// Execute the request, expecting no response body.
	_, err = s.client.Do(req, nil)
	return err
}
