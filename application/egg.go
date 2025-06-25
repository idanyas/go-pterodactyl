package application

import (
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type eggsService struct {
	client requester.Requester
	nestID int
}

func NewEggsService(client requester.Requester, nestID int) *eggsService {
	return &eggsService{client: client, nestID: nestID}
}

func (s *eggsService) List(options *api.PaginationOptions) ([]*api.Egg, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d/eggs", s.nestID)
	req, err := s.client.NewRequest("GET", endpoint, nil, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create egg list request: %w", err)
	}

	response := &api.PaginatedResponse[api.Egg]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Egg, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *eggsService) ListAll() ([]*api.Egg, error) {
	allEggs := make([]*api.Egg, 0, 100)
	options := &api.PaginationOptions{PerPage: 100, Page: 1}

	for {
		eggs, meta, err := s.List(options)
		if err != nil {
			return nil, err
		}

		allEggs = append(allEggs, eggs...)

		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}

		options.Page++
	}

	return allEggs, nil
}

func (s *eggsService) Get(eggID int) (*api.Egg, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d/eggs/%d", s.nestID, eggID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get egg request: %w", err)
	}

	response := &api.ListItem[api.Egg]{}
	_, err = s.client.Do(req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}
