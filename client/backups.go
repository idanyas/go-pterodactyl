package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// CreateBackupRequest defines the request body for creating a new backup.
type CreateBackupRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=1"`
	Ignored  string `json:"ignored,omitempty"`
	IsLocked bool   `json:"is_locked,omitempty"`
}

// ListBackups retrieves all backups for a server.
func (c *client) ListBackups(ctx context.Context, serverID string) ([]*models.Backup, error) {
	path := fmt.Sprintf("client/servers/%s/backups", serverID)
	var response struct {
		Data []struct {
			Attributes models.Backup `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	backups := make([]*models.Backup, len(response.Data))
	for i, item := range response.Data {
		backups[i] = &item.Attributes
	}
	return backups, nil
}

// GetBackup retrieves details for a specific backup.
func (c *client) GetBackup(ctx context.Context, serverID, backupUUID string) (*models.Backup, error) {
	path := fmt.Sprintf("client/servers/%s/backups/%s", serverID, backupUUID)
	var response struct {
		Attributes models.Backup `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// CreateBackup creates a new backup for the server.
func (c *client) CreateBackup(ctx context.Context, serverID string, req CreateBackupRequest) (*models.Backup, error) {
	path := fmt.Sprintf("client/servers/%s/backups", serverID)
	var response struct {
		Attributes models.Backup `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// GetBackupDownloadURL retrieves a pre-signed URL to download a backup.
func (c *client) GetBackupDownloadURL(ctx context.Context, serverID, backupUUID string) (*models.SignedURL, error) {
	path := fmt.Sprintf("client/servers/%s/backups/%s/download", serverID, backupUUID)
	var response struct {
		Attributes models.SignedURL `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// RestoreBackup restores a server from a backup.
func (c *client) RestoreBackup(ctx context.Context, serverID, backupUUID string, truncate bool) error {
	path := fmt.Sprintf("client/servers/%s/backups/%s/restore", serverID, backupUUID)
	req := map[string]bool{"truncate": truncate}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// ToggleBackupLock toggles the lock status of a backup.
func (c *client) ToggleBackupLock(ctx context.Context, serverID, backupUUID string) error {
	path := fmt.Sprintf("client/servers/%s/backups/%s/lock", serverID, backupUUID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	return err
}

// DeleteBackup permanently deletes a backup.
func (c *client) DeleteBackup(ctx context.Context, serverID, backupUUID string) error {
	path := fmt.Sprintf("client/servers/%s/backups/%s", serverID, backupUUID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
