package client

import (
	"context"
	"fmt"
	"net/http"
)

// RenameServer updates the name and description of a server.
func (c *client) RenameServer(ctx context.Context, serverID, name, description string) error {
	path := fmt.Sprintf("client/servers/%s/settings/rename", serverID)
	req := map[string]string{"name": name, "description": description}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// ReinstallServer triggers a reinstallation of the server.
func (c *client) ReinstallServer(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("client/servers/%s/settings/reinstall", serverID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	return err
}

// UpdateDockerImage changes the Docker image used by the server.
func (c *client) UpdateDockerImage(ctx context.Context, serverID, dockerImage string) error {
	path := fmt.Sprintf("client/servers/%s/settings/docker-image", serverID)
	req := map[string]string{"docker_image": dockerImage}
	_, err := c.client.Do(ctx, http.MethodPut, path, req, nil)
	return err
}
