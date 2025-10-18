package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// CreateServerRequest defines the request body for creating a new server.
type CreateServerRequest struct {
	Name              string                 `json:"name" validate:"required,min=1"`
	User              int                    `json:"user" validate:"required,gt=0"`
	Egg               int                    `json:"egg" validate:"required,gt=0"`
	DockerImage       string                 `json:"docker_image,omitempty"`
	Startup           string                 `json:"startup,omitempty"`
	Environment       map[string]string      `json:"environment,omitempty"`
	Limits            models.Limits          `json:"limits" validate:"required"`
	FeatureLimits     models.FeatureLimits   `json:"feature_limits" validate:"required"`
	Allocation        CreateServerAllocation `json:"allocation" validate:"required"`
	Deploy            *CreateServerDeploy    `json:"deploy,omitempty"`
	Description       string                 `json:"description,omitempty"`
	ExternalID        string                 `json:"external_id,omitempty"`
	SkipScripts       bool                   `json:"skip_scripts,omitempty"`
	StartOnCompletion bool                   `json:"start_on_completion,omitempty"`
}

// CreateServerAllocation defines the allocation configuration for a new server.
type CreateServerAllocation struct {
	Default    int   `json:"default" validate:"required,gt=0"`
	Additional []int `json:"additional,omitempty"`
}

// CreateServerDeploy defines the deployment configuration for a new server.
type CreateServerDeploy struct {
	Locations   []int    `json:"locations" validate:"required,min=1,dive,gt=0"`
	DedicatedIP bool     `json:"dedicated_ip,omitempty"`
	PortRange   []string `json:"port_range,omitempty"`
}

// UpdateServerDetailsRequest defines the request body for updating a server's details.
type UpdateServerDetailsRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=1"`
	User        int    `json:"user,omitempty" validate:"omitempty,gt=0"`
	ExternalID  string `json:"external_id,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateServerBuildRequest defines the request body for updating a server's build configuration.
type UpdateServerBuildRequest struct {
	AllocationID      int                  `json:"allocation_id" validate:"required,gt=0"`
	AddAllocations    []int                `json:"add_allocations,omitempty"`
	RemoveAllocations []int                `json:"remove_allocations,omitempty"`
	Memory            int64                `json:"memory" validate:"required,gt=0"`
	Swap              int64                `json:"swap" validate:"gte=0"`
	Disk              int64                `json:"disk" validate:"required,gt=0"`
	IO                int64                `json:"io" validate:"gte=10,lte=1000"`
	CPU               int64                `json:"cpu" validate:"gte=0"`
	Threads           string               `json:"threads,omitempty"`
	FeatureLimits     models.FeatureLimits `json:"feature_limits" validate:"required"`
}

// UpdateServerStartupRequest defines the request body for updating a server's startup configuration.
type UpdateServerStartupRequest struct {
	Startup     string            `json:"startup" validate:"required"`
	Environment map[string]string `json:"environment" validate:"required"`
	Egg         int               `json:"egg" validate:"required,gt=0"`
	Image       string            `json:"image" validate:"required"`
	SkipScripts bool              `json:"skip_scripts"`
}

// ListServers retrieves a paginated list of all servers.
func (c *client) ListServers(ctx context.Context, options pagination.ListOptions) ([]*models.Server, *pagination.Paginator[*models.Server], error) {
	return pagination.New[*models.Server](ctx, c.client, "application/servers", options)
}

// GetServer retrieves details for a specific server by its internal ID.
func (c *client) GetServer(ctx context.Context, id int) (*models.Server, error) {
	if id <= 0 {
		return nil, fmt.Errorf("server ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/servers/%d", id)
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get server %d: %w", id, err)
	}
	return &response.Attributes, nil
}

// GetServerExternal retrieves details for a specific server by its external ID.
func (c *client) GetServerExternal(ctx context.Context, externalID string) (*models.Server, error) {
	if externalID == "" {
		return nil, fmt.Errorf("external ID cannot be empty")
	}

	path := fmt.Sprintf("application/servers/external/%s", externalID)
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get server with external ID %s: %w", externalID, err)
	}
	return &response.Attributes, nil
}

// CreateServer creates a new server.
func (c *client) CreateServer(ctx context.Context, req CreateServerRequest) (*models.Server, error) {
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "application/servers", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	return &response.Attributes, nil
}

// UpdateServerDetails updates a server's name, owner, and description.
func (c *client) UpdateServerDetails(ctx context.Context, serverID int, req UpdateServerDetailsRequest) (*models.Server, error) {
	if serverID <= 0 {
		return nil, fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d/details", serverID)
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update server %d details: %w", serverID, err)
	}
	return &response.Attributes, nil
}

// UpdateServerBuild updates a server's resource limits and allocations.
func (c *client) UpdateServerBuild(ctx context.Context, serverID int, req UpdateServerBuildRequest) (*models.Server, error) {
	if serverID <= 0 {
		return nil, fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d/build", serverID)
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update server %d build: %w", serverID, err)
	}
	return &response.Attributes, nil
}

// UpdateServerStartup updates a server's Egg, startup command, and environment variables.
func (c *client) UpdateServerStartup(ctx context.Context, serverID int, req UpdateServerStartupRequest) (*models.Server, error) {
	if serverID <= 0 {
		return nil, fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d/startup", serverID)
	var response struct {
		Attributes models.Server `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update server %d startup: %w", serverID, err)
	}
	return &response.Attributes, nil
}

// SuspendServer suspends a server.
func (c *client) SuspendServer(ctx context.Context, serverID int) error {
	if serverID <= 0 {
		return fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d/suspend", serverID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to suspend server %d: %w", serverID, err)
	}
	return nil
}

// UnsuspendServer unsuspends a server.
func (c *client) UnsuspendServer(ctx context.Context, serverID int) error {
	if serverID <= 0 {
		return fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d/unsuspend", serverID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to unsuspend server %d: %w", serverID, err)
	}
	return nil
}

// ReinstallServer triggers a reinstallation for a server.
func (c *client) ReinstallServer(ctx context.Context, serverID int) error {
	if serverID <= 0 {
		return fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d/reinstall", serverID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to reinstall server %d: %w", serverID, err)
	}
	return nil
}

// DeleteServer permanently deletes a server.
func (c *client) DeleteServer(ctx context.Context, serverID int, force bool) error {
	if serverID <= 0 {
		return fmt.Errorf("server ID must be positive, got %d", serverID)
	}

	path := fmt.Sprintf("application/servers/%d", serverID)
	if force {
		path = fmt.Sprintf("application/servers/%d/force", serverID)
	}
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete server %d: %w", serverID, err)
	}
	return nil
}
