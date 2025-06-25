package application

import (
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type nestsService struct {
	client requester.Requester
}

func (s *nestsService) Eggs(nestID int) EggsService {
	return NewEggsService(s.client, nestID)
}

func NewNestsService(client requester.Requester) *nestsService {
	return &nestsService{client: client}
}

func (s *nestsService) List(options *api.PaginationOptions) ([]*api.Nest, *api.Meta, error) {
	req, err := s.client.NewRequest("GET", "/api/application/nests", nil, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create nest list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Nest]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Nest, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *nestsService) ListAll() ([]*api.Nest, error) {
	allNests := make([]*api.Nest, 0, 100)
	options := &api.PaginationOptions{PerPage: 100, Page: 1}

	for {
		nests, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}

		allNests = append(allNests, nests...)

		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}

		options.Page++
	}

	return allNests, nil
}

// Get fetches a single nest by its ID.
func (s *nestsService) Get(nestID int) (*api.Nest, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d", nestID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get nest request: %w", err)
	}

	response := &api.ListItem[api.Nest]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}
