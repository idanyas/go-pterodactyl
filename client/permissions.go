package client

import (
	"context"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// GetSystemPermissions retrieves all available permissions that can be assigned to subusers.
func (c *client) GetSystemPermissions(ctx context.Context) (*models.SystemPermissions, error) {
	var response struct {
		Attributes models.SystemPermissions `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, "client/permissions", nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}
