package clientapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type networkService struct {
	client           requester.Requester
	serverIdentifier string
}

// newNetworkService creates a new network service.
func newNetworkService(client requester.Requester, serverIdentifier string) *networkService {
	return &networkService{client: client, serverIdentifier: serverIdentifier}
}

// ListAllocations retrieves all network allocations for the server.
func (s *networkService) ListAllocations(options api.PaginationOptions) ([]*api.Allocation, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations", s.serverIdentifier)
	req, err := s.client.NewRequest("GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list allocations request: %w", err)
	}

	res := &api.PaginatedResponse[api.Allocation]{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Allocation, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

// AssignAllocation requests that a new allocation be automatically assigned to the server.
func (s *networkService) AssignAllocation() (*api.Allocation, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations", s.serverIdentifier)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create assign allocation request: %w", err)
	}

	res := &api.AllocationResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// SetAllocationNote updates the notes for a specific allocation.
func (s *networkService) SetAllocationNote(allocationID int, options api.AllocationNoteOptions) (*api.Allocation, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal allocation note options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d", s.serverIdentifier, allocationID)
	req, err := s.client.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create set allocation note request: %w", err)
	}

	res := &api.AllocationResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// SetPrimaryAllocation designates an allocation as the primary one for the server.
func (s *networkService) SetPrimaryAllocation(allocationID int) (*api.Allocation, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d/primary", s.serverIdentifier, allocationID)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create set primary allocation request: %w", err)
	}

	res := &api.AllocationResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// UnassignAllocation removes a network allocation from the server.
func (s *networkService) UnassignAllocation(allocationID int) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d", s.serverIdentifier, allocationID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create unassign allocation request: %w", err)
	}

	_, err = s.client.Do(req, nil)
	return err
}
