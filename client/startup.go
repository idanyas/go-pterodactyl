package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// GetStartupConfig retrieves the startup configuration and environment variables for a server.
func (c *client) GetStartupConfig(ctx context.Context, serverID string) (*models.StartupConfiguration, error) {
	path := fmt.Sprintf("client/servers/%s/startup", serverID)
	var response struct {
		Data []struct {
			Attributes models.StartupVariable `json:"attributes"`
		} `json:"data"`
		Meta struct {
			StartupCommand    string            `json:"startup_command"`
			RawStartupCommand string            `json:"raw_startup_command"`
			DockerImages      map[string]string `json:"docker_images"`
		} `json:"meta"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	vars := make([]*models.StartupVariable, len(response.Data))
	for i, item := range response.Data {
		vars[i] = &item.Attributes
	}

	config := &models.StartupConfiguration{
		Variables:         vars,
		StartupCommand:    response.Meta.StartupCommand,
		RawStartupCommand: response.Meta.RawStartupCommand,
		DockerImages:      response.Meta.DockerImages,
	}
	return config, nil
}

// UpdateStartupVariable updates a specific startup environment variable.
func (c *client) UpdateStartupVariable(ctx context.Context, serverID, key, value string) (*models.StartupVariable, error) {
	path := fmt.Sprintf("client/servers/%s/startup/variable", serverID)
	req := map[string]string{"key": key, "value": value}
	var response struct {
		Attributes models.StartupVariable `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPut, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}
