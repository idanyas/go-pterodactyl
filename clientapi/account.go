package clientapi

import (
	"bytes"
	"encoding/json"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type accountService struct{ client requester.Requester }

func newAccountService(client requester.Requester) AccountService {
	return &accountService{client: client}
}

func (s *accountService) GetDetails() (*api.Account, error) {
	req, err := s.client.NewRequest("GET", "/api/client/account", nil, nil)
	if err != nil {
		return nil, err
	}

	res := &api.ListItem[api.Account]{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

func (s *accountService) GetTwoFactorDetails() (*api.TwoFactorDetails, error) {
	req, err := s.client.NewRequest("GET", "/api/client/account/two-factor", nil, nil)
	if err != nil {
		return nil, err
	}

	res := &api.TwoFactorDetails{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *accountService) EnableTwoFactor(options api.TwoFactorEnableOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest("POST", "/api/client/account/two-factor", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *accountService) DisableTwoFactor(options api.TwoFactorDisableOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest("DELETE", "/api/client/account/two-factor", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *accountService) UpdateEmail(options api.UpdateEmailOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest("PUT", "/api/client/account/email", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *accountService) UpdatePassword(options api.UpdatePasswordOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest("PUT", "/api/client/account/password", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *accountService) APIKeys() APIKeysService {
	return newAPIKeysService(s.client)
}
