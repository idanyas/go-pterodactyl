package clientapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type apiKeysService struct{ client requester.Requester }

func newAPIKeysService(client requester.Requester) APIKeysService {
	return &apiKeysService{client: client}
}

func (s *apiKeysService) List(options api.PaginationOptions) ([]*api.APIKey, *api.Meta, error) {
	req, err := s.client.NewRequest("GET", "/api/client/account/api-keys", nil, &options)
	if err != nil {
		return nil, nil, err
	}

	res := &api.PaginatedResponse[api.APIKey]{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.APIKey, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

func (s *apiKeysService) Create(options api.APIKeyCreateOptions) (*api.APIKey, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", "/api/client/account/api-keys", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}

	res := &api.APIKeyCreateResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	apiKey := &res.Attributes

	apiKey.Token = &res.Meta.SecretToken

	return apiKey, nil
}

func (s *apiKeysService) Delete(identifier string) error {
	endpoint := fmt.Sprintf("/api/client/account/api-keys/%s", identifier)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
