package clientapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type databasesService struct {
	client           requester.Requester
	serverIdentifier string
}

func newDatabasesService(client requester.Requester, serverIdentifier string) *databasesService {
	return &databasesService{client: client, serverIdentifier: serverIdentifier}
}

func (s *databasesService) List(options api.PaginationOptions) ([]*api.ClientDatabase, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/databases", s.serverIdentifier)
	req, err := s.client.NewRequest("GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list databases request: %w", err)
	}

	res := &api.PaginatedResponse[api.ClientDatabase]{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.ClientDatabase, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

func processCreateOrRotateResponse(res *api.ClientDatabaseCreateResponse) *api.ClientDatabase {
	db := res.Attributes
	if res.Relationships != nil && res.Relationships.Password != nil && res.Relationships.Password.Attributes != nil {
		db.Password = res.Relationships.Password.Attributes.Password
	}
	return db
}

func (s *databasesService) Create(options api.ClientDatabaseCreateOptions) (*api.ClientDatabase, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create database options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/databases", s.serverIdentifier)
	req, err := s.client.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new database request: %w", err)
	}

	// Use the special response struct to capture the password from relationships.
	res := &api.ClientDatabaseCreateResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return processCreateOrRotateResponse(res), nil
}

func (s *databasesService) RotatePassword(databaseID string) (*api.ClientDatabase, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/databases/%s/rotate-password", s.serverIdentifier, databaseID)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create rotate password request: %w", err)
	}

	// This endpoint also returns the special response with the password.
	res := &api.ClientDatabaseCreateResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return processCreateOrRotateResponse(res), nil
}

func (s *databasesService) Delete(databaseID string) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/databases/%s", s.serverIdentifier, databaseID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete database request: %w", err)
	}

	_, err = s.client.Do(req, nil)
	return err
}
