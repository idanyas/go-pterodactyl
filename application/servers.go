package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
	"io"
	"net/url"
)

type serversService struct {
	client requester.Requester
}

// NewServersService is the exported constructor.
func NewServersService(client requester.Requester) ServersService {
	return &serversService{client: client}
}

func (s *serversService) List(options api.PaginationOptions) ([]*api.Server, *api.Meta, error) {
	req, err := s.client.NewRequest("GET", "/api/application/servers", nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create server list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Server, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *serversService) ListAll() ([]*api.Server, error) {
	allServers := make([]*api.Server, 0, 100)
	options := api.PaginationOptions{PerPage: 100, Page: 1}

	for {
		servers, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}
		allServers = append(allServers, servers...)
		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}
		options.Page++
	}
	return allServers, nil
}

func (s *serversService) Get(id int) (*api.Server, error) {
	endpoint := fmt.Sprintf("/api/application/servers/%d", id)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get server request: %w", err)
	}

	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) GetExternal(externalID string) (*api.Server, error) {
	endpoint := fmt.Sprintf("/api/application/servers/external/%s", url.PathEscape(externalID))
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get external server request: %w", err)
	}

	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) Create(options api.ServerCreateOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create server options: %w", err)
	}
	req, err := s.client.NewRequest("POST", "/api/application/servers", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) UpdateDetails(serverID int, options api.ServerUpdateDetailsOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d/details", serverID)
	req, err := s.client.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) UpdateBuild(serverID int, options api.ServerUpdateBuildOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d/build", serverID)
	req, err := s.client.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) UpdateStartup(serverID int, options api.ServerUpdateStartupOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d/startup", serverID)
	req, err := s.client.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) Suspend(serverID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/suspend", serverID)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *serversService) Unsuspend(serverID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/unsuspend", serverID)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *serversService) Reinstall(serverID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/reinstall", serverID)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *serversService) Delete(serverID int, force bool) error {
	var body io.Reader
	if force {
		jsonBytes, err := json.Marshal(api.ServerDeleteOptions{Force: true})
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBytes)
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d", serverID)
	req, err := s.client.NewRequest("DELETE", endpoint, body, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Databases returns a specialized service for managing databases for a specific server.
func (s *serversService) Databases(serverID int) DatabaseService {
	return newDatabaseService(s.client, serverID)
}
