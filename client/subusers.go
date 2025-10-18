package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// ListSubusers retrieves all users with access to a server.
func (c *client) ListSubusers(ctx context.Context, serverID string) ([]*models.Subuser, error) {
	path := fmt.Sprintf("client/servers/%s/users", serverID)
	var response struct {
		Data []struct {
			Attributes models.Subuser `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	subusers := make([]*models.Subuser, len(response.Data))
	for i, item := range response.Data {
		subusers[i] = &item.Attributes
	}
	return subusers, nil
}

// GetSubuser retrieves details for a specific subuser.
func (c *client) GetSubuser(ctx context.Context, serverID, userUUID string) (*models.Subuser, error) {
	path := fmt.Sprintf("client/servers/%s/users/%s", serverID, userUUID)
	var response struct {
		Attributes models.Subuser `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// CreateSubuser invites a new user to the server with specific permissions.
func (c *client) CreateSubuser(ctx context.Context, serverID, email string, permissions []string) (*models.Subuser, error) {
	path := fmt.Sprintf("client/servers/%s/users", serverID)
	req := map[string]interface{}{
		"email":       email,
		"permissions": permissions,
	}
	var response struct {
		Attributes models.Subuser `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// UpdateSubuser updates the permissions for an existing subuser.
func (c *client) UpdateSubuser(ctx context.Context, serverID, userUUID string, permissions []string) (*models.Subuser, error) {
	path := fmt.Sprintf("client/servers/%s/users/%s", serverID, userUUID)
	req := map[string]interface{}{"permissions": permissions}
	var response struct {
		Attributes models.Subuser `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DeleteSubuser removes a user's access from the server.
func (c *client) DeleteSubuser(ctx context.Context, serverID, userUUID string) error {
	path := fmt.Sprintf("client/servers/%s/users/%s", serverID, userUUID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
