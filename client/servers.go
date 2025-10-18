package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// ValidPowerSignals are the valid power action signals.
var ValidPowerSignals = []string{"start", "stop", "restart", "kill"}

// ListServers retrieves a paginated list of all servers the user has access to.
func (c *client) ListServers(ctx context.Context, options pagination.ListOptions) ([]*models.Server, *pagination.Paginator[*models.Server], error) {
	return pagination.New[*models.Server](ctx, c.client, "client", options)
}

// GetServer retrieves details for a specific server.
func (c *client) GetServer(ctx context.Context, serverID string) (*models.Server, error) {
	if serverID == "" {
		return nil, fmt.Errorf("server ID cannot be empty")
	}

	path := fmt.Sprintf("client/servers/%s", serverID)
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get server %s: %w", serverID, err)
	}
	return &response.Attributes, nil
}

// GetServerResources retrieves real-time resource usage for a server.
func (c *client) GetServerResources(ctx context.Context, serverID string) (*models.Stats, error) {
	if serverID == "" {
		return nil, fmt.Errorf("server ID cannot be empty")
	}

	path := fmt.Sprintf("client/servers/%s/resources", serverID)
	var response struct {
		Attributes models.Stats `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get server %s resources: %w", serverID, err)
	}
	return &response.Attributes, nil
}

// SendPowerAction sends a power signal to a server (start, stop, restart, kill).
func (c *client) SendPowerAction(ctx context.Context, serverID, signal string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}
	if signal == "" {
		return fmt.Errorf("signal cannot be empty")
	}

	// Validate signal
	validSignal := false
	for _, valid := range ValidPowerSignals {
		if signal == valid {
			validSignal = true
			break
		}
	}
	if !validSignal {
		return fmt.Errorf("invalid power signal %q, must be one of: %v", signal, ValidPowerSignals)
	}

	path := fmt.Sprintf("client/servers/%s/power", serverID)
	req := map[string]string{"signal": signal}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("failed to send power action %s to server %s: %w", signal, serverID, err)
	}
	return nil
}

// SendCommand executes a command in the server console.
func (c *client) SendCommand(ctx context.Context, serverID, command string) error {
	if serverID == "" {
		return fmt.Errorf("server ID cannot be empty")
	}
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	path := fmt.Sprintf("client/servers/%s/command", serverID)
	req := map[string]string{"command": command}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	if err != nil {
		return fmt.Errorf("failed to send command to server %s: %w", serverID, err)
	}
	return nil
}
