package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type LocationService struct {
	client requester.Requester
}

func NewLocationService(client requester.Requester) *LocationService {
	return &LocationService{client: client}
}

func (s *LocationService) List(options *api.PaginationOptions) ([]*api.Location, *api.Meta, error) {
	endpoint := "/api/application/locations"
	req, err := s.client.NewRequest("GET", endpoint, nil, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create location list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Location]{}

	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Location, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}

	return results, &response.Meta, nil
}

func (s *LocationService) ListAll() ([]*api.Location, error) {
	allLocations := make([]*api.Location, 0, 100)
	options := &api.PaginationOptions{PerPage: 100, Page: 1}
	for {
		locations, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}

		allLocations = append(allLocations, locations...)

		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}

		options.Page++
	}

	return allLocations, nil

}

func (s *LocationService) Get(id int) (*api.Location, error) {
	endpoint := fmt.Sprintf("/api/application/locations/%d", id)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create location get request: %w", err)
	}
	response := &api.ListItem[api.Location]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *LocationService) Create(options api.LocationCreateOptions) (*api.Location, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create location options: %w", err)
	}

	req, err := s.client.NewRequest("POST", "/api/application/locations", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new location request: %w", err)
	}

	response := &api.ListItem[api.Location]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *LocationService) Update(id int, options api.LocationUpdateOptions) (*api.Location, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update location options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/locations/%d", id)
	req, err := s.client.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update location request: %w", err)
	}

	response := &api.ListItem[api.Location]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *LocationService) Delete(id int) error {
	endpoint := fmt.Sprintf("/api/application/locations/%d", id)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete location request: %w", err)
	}

	_, err = s.client.Do(req, nil)
	return err
}
