package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// ListDatabases retrieves all databases associated with a server.
func (c *client) ListDatabases(ctx context.Context, serverID string) ([]*models.Database, error) {
	path := fmt.Sprintf("client/servers/%s/databases", serverID)
	var response struct {
		Data []struct {
			Attributes models.Database `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	databases := make([]*models.Database, len(response.Data))
	for i, item := range response.Data {
		databases[i] = &item.Attributes
	}
	return databases, nil
}

// CreateDatabase creates a new database for the server.
func (c *client) CreateDatabase(ctx context.Context, serverID, database, remote string) (*models.Database, error) {
	path := fmt.Sprintf("client/servers/%s/databases", serverID)
	req := map[string]string{"database": database, "remote": remote}
	var response struct {
		Attributes models.Database `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// RotateDatabasePassword generates a new password for an existing database.
func (c *client) RotateDatabasePassword(ctx context.Context, serverID, databaseID string) (*models.Database, error) {
	path := fmt.Sprintf("client/servers/%s/databases/%s/rotate-password", serverID, databaseID)
	var response struct {
		Attributes models.Database `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DeleteDatabase permanently deletes a database.
func (c *client) DeleteDatabase(ctx context.Context, serverID, databaseID string) error {
	path := fmt.Sprintf("client/servers/%s/databases/%s", serverID, databaseID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
