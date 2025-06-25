package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type nodesService struct {
	client requester.Requester
}

func NewNodesService(client requester.Requester) NodesService {
	return &nodesService{client: client}
}

func (s *nodesService) List(options *api.PaginationOptions) ([]*api.Node, *api.Meta, error) {
	req, err := s.client.NewRequest("GET", "/api/application/nodes", nil, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create node list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Node]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Node, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *nodesService) ListAll() ([]*api.Node, error) {
	allNodes := make([]*api.Node, 0, 100)
	options := &api.PaginationOptions{PerPage: 100, Page: 1}

	for {
		nodes, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}

		allNodes = append(allNodes, nodes...)

		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}

		options.Page++
	}

	return allNodes, nil
}

func (s *nodesService) Get(id int) (*api.Node, error) {
	endpoint := fmt.Sprintf("/api/application/nodes/%d", id)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get node request: %w", err)
	}

	// A successful response for a single node is wrapped in a ListItem.
	response := &api.ListItem[api.Node]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *nodesService) GetConfiguration(nodeID int) (*api.NodeConfiguration, error) {
	endpoint := fmt.Sprintf("/api/application/nodes/%d/configuration", nodeID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get node configuration request: %w", err)
	}

	// The response for this endpoint is a unique structure, not a ListItem.
	// Our generic `Do` method handles this perfectly.
	response := &api.NodeConfiguration{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *nodesService) Create(options api.NodeCreateOptions) (*api.Node, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create node options: %w", err)
	}

	req, err := s.client.NewRequest("POST", "/api/application/nodes", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new node request: %w", err)
	}

	response := &api.ListItem[api.Node]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *nodesService) Update(id int, options api.NodeUpdateOptions) (*api.Node, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update node options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/nodes/%d", id)
	req, err := s.client.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update node request: %w", err)
	}

	response := &api.ListItem[api.Node]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *nodesService) Delete(id int) error {
	endpoint := fmt.Sprintf("/api/application/nodes/%d", id)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete node request: %w", err)
	}

	_, err = s.client.Do(req, nil)
	return err
}

func (s *nodesService) Allocations(nodeID int) AllocationsService {
	return newAllocationsService(s.client, nodeID)
}
