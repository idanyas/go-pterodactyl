package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type databaseService struct {
	client   requester.Requester
	serverID int
}

func newDatabaseService(client requester.Requester, serverID int) *databaseService {
	return &databaseService{client: client, serverID: serverID}
}

func (s *databaseService) List(options api.PaginationOptions) ([]*api.Database, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases", s.serverID)
	req, err := s.client.NewRequest("GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Database]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Database, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

// Get fetches a single database by its ID for the configured server.
func (s *databaseService) Get(databaseID int) (*api.Database, error) {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases/%d", s.serverID, databaseID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get database request: %w", err)
	}

	response := &api.ListItem[api.Database]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

// Create creates a new database for the configured server.
// The response includes the new password, so we need to return the full Database object.
func (s *databaseService) Create(options api.DatabaseCreateOptions) (*api.Database, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create database options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/servers/%d/databases", s.serverID)
	req, err := s.client.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new database request: %w", err)
	}

	response := &api.ListItem[api.Database]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	// The response for this endpoint includes a `relationships.password` field
	// which is not part of the standard `Database` struct. We might need a custom struct
	// if we want to capture that password. For now, returning the created DB object is fine.
	return response.Attributes, nil
}

// ResetPassword requests a password reset for a specific database.
func (s *databaseService) ResetPassword(databaseID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases/%d/reset-password", s.serverID, databaseID)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create reset database password request: %w", err)
	}

	// This endpoint returns 204 No Content.
	_, err = s.client.Do(req, nil)
	return err
}

// Delete deletes a specific database from the configured server.
func (s *databaseService) Delete(databaseID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases/%d", s.serverID, databaseID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete database request: %w", err)
	}

	// This endpoint returns 204 No Content.
	_, err = s.client.Do(req, nil)
	return err
}
