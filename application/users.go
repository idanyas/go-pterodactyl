package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type usersService struct {
	client requester.Requester
}

func NewUsersService(client requester.Requester) *usersService {
	return &usersService{client: client}
}

func (s *usersService) List(options *api.PaginationOptions) ([]*api.User, *api.Meta, error) {
	req, err := s.client.NewRequest("GET", "/api/application/users", nil, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user list request: %w", err)
	}

	response := &api.PaginatedResponse[api.User]{}

	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.User, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}

	return results, &response.Meta, nil
}

func (s *usersService) ListAll() ([]*api.User, error) {
	allUsers := make([]*api.User, 0, 100)
	options := &api.PaginationOptions{PerPage: 100, Page: 1}

	for {
		users, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, users...)

		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}

		options.Page++
	}

	return allUsers, nil
}

func (s *usersService) Get(id int) (*api.User, error) {
	// Construct the correct endpoint URL with the user's ID.
	endpoint := fmt.Sprintf("/api/application/users/%d", id)

	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user get request: %w", err)
	}

	response := &api.ListItem[api.User]{}

	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *usersService) GetExternalID(externalId string) (*api.User, error) {
	endpoint := fmt.Sprintf("/api/application/users/external/%s", externalId)

	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user get request: %w", err)
	}
	response := &api.ListItem[api.User]{}

	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil

}

func (s *usersService) Create(options api.UserCreateOptions) (*api.User, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create user options: %w", err)
	}

	req, err := s.client.NewRequest("POST", "/api/application/users", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user request: %w", err)
	}

	response := &api.ListItem[api.User]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *usersService) Update(id int, options api.UserUpdateOptions) (*api.User, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update user options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/users/%d", id)
	req, err := s.client.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update user request: %w", err)
	}

	response := &api.ListItem[api.User]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *usersService) Delete(id int) error {
	endpoint := fmt.Sprintf("/api/application/users/%d", id)

	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user request: %w", err)
	}

	_, err = s.client.Do(req, nil)

	return err
}
