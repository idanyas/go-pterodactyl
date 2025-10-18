package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// CreateServerDatabaseRequest defines the request body for creating a server database.
type CreateServerDatabaseRequest struct {
	Database string `json:"database" validate:"required,min=1"`
	Remote   string `json:"remote" validate:"required"`
	Host     int    `json:"host" validate:"required,gt=0"`
}

// UpdateServerDatabaseRequest defines the request body for updating a server database.
type UpdateServerDatabaseRequest struct {
	Remote string `json:"remote,omitempty" validate:"omitempty,min=1"`
}

// ListServerDatabases retrieves all databases for a specific server.
func (c *client) ListServerDatabases(ctx context.Context, serverID int, options pagination.ListOptions) ([]*models.ApplicationDatabase, *pagination.Paginator[*models.ApplicationDatabase], error) {
	path := fmt.Sprintf("application/servers/%d/databases", serverID)
	return pagination.New[*models.ApplicationDatabase](ctx, c.client, path, options)
}

// GetServerDatabase retrieves details for a specific database.
func (c *client) GetServerDatabase(ctx context.Context, serverID, databaseID int) (*models.ApplicationDatabase, error) {
	path := fmt.Sprintf("application/servers/%d/databases/%d", serverID, databaseID)
	var response struct {
		Attributes models.ApplicationDatabase `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// CreateServerDatabase creates a new database for a server.
func (c *client) CreateServerDatabase(ctx context.Context, serverID int, req CreateServerDatabaseRequest) (*models.ApplicationDatabase, error) {
	path := fmt.Sprintf("application/servers/%d/databases", serverID)
	var response struct {
		Attributes models.ApplicationDatabase `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// UpdateServerDatabase updates a server database configuration.
func (c *client) UpdateServerDatabase(ctx context.Context, serverID, databaseID int, req UpdateServerDatabaseRequest) (*models.ApplicationDatabase, error) {
	path := fmt.Sprintf("application/servers/%d/databases/%d", serverID, databaseID)
	var response struct {
		Attributes models.ApplicationDatabase `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// ResetServerDatabasePassword generates a new password for a server database.
func (c *client) ResetServerDatabasePassword(ctx context.Context, serverID, databaseID int) error {
	path := fmt.Sprintf("application/servers/%d/databases/%d/reset-password", serverID, databaseID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	return err
}

// DeleteServerDatabase permanently deletes a server database.
func (c *client) DeleteServerDatabase(ctx context.Context, serverID, databaseID int) error {
	path := fmt.Sprintf("application/servers/%d/databases/%d", serverID, databaseID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
