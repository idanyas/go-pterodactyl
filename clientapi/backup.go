package clientapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/internal/requester"
)

type backupsService struct {
	client           requester.Requester
	serverIdentifier string
}

// newBackupsService creates a new backups service.
func newBackupsService(client requester.Requester, serverIdentifier string) *backupsService {
	return &backupsService{client: client, serverIdentifier: serverIdentifier}
}

// List retrieves all backups for the server.
func (s *backupsService) List(options api.PaginationOptions) ([]*api.Backup, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups", s.serverIdentifier)
	req, err := s.client.NewRequest("GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list backups request: %w", err)
	}

	res := &api.PaginatedResponse[api.Backup]{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Backup, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

// Create sends a request to begin a new backup creation process.
func (s *backupsService) Create(options api.BackupCreateOptions) (*api.Backup, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create backup options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/backups", s.serverIdentifier)
	req, err := s.client.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new backup request: %w", err)
	}

	res := &api.BackupResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// Details retrieves the details of a specific backup by its UUID.
func (s *backupsService) Details(uuid string) (*api.Backup, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup details request: %w", err)
	}

	res := &api.BackupResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

func (s *backupsService) Download(uuid string) (*api.BackupDownload, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s/download", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup download request: %w", err)
	}

	res := &api.BackupDownloadResponse{}
	_, err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

func (s *backupsService) Delete(uuid string) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete backup request: %w", err)
	}

	_, err = s.client.Do(req, nil)
	return err
}
