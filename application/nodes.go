package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// CreateNodeRequest defines the request body for creating a new node.
type CreateNodeRequest struct {
	Name               string `json:"name" validate:"required,min=1"`
	LocationID         int    `json:"location_id" validate:"required,gt=0"`
	FQDN               string `json:"fqdn" validate:"required,fqdn"`
	Scheme             string `json:"scheme" validate:"required,oneof=http https"`
	Memory             int64  `json:"memory" validate:"required,gt=0"`
	MemoryOverallocate int64  `json:"memory_overallocate" validate:"gte=-1"`
	Disk               int64  `json:"disk" validate:"required,gt=0"`
	DiskOverallocate   int64  `json:"disk_overallocate" validate:"gte=-1"`
	UploadSize         int64  `json:"upload_size" validate:"required,gt=0"`
	DaemonSFTP         int    `json:"daemon_sftp" validate:"required,gt=0,lt=65536"`
	DaemonListen       int    `json:"daemon_listen" validate:"required,gt=0,lt=65536"`
	DaemonBase         string `json:"daemon_base" validate:"required"`
	Public             bool   `json:"public"`
	BehindProxy        bool   `json:"behind_proxy"`
	MaintenanceMode    bool   `json:"maintenance_mode"`
	Description        string `json:"description,omitempty"`
}

// UpdateNodeRequest defines the request body for updating a node.
type UpdateNodeRequest struct {
	Name               string `json:"name,omitempty" validate:"omitempty,min=1"`
	LocationID         int    `json:"location_id,omitempty" validate:"omitempty,gt=0"`
	FQDN               string `json:"fqdn,omitempty" validate:"omitempty,fqdn"`
	Scheme             string `json:"scheme,omitempty" validate:"omitempty,oneof=http https"`
	Memory             int64  `json:"memory,omitempty" validate:"omitempty,gt=0"`
	MemoryOverallocate int64  `json:"memory_overallocate,omitempty" validate:"omitempty,gte=-1"`
	Disk               int64  `json:"disk,omitempty" validate:"omitempty,gt=0"`
	DiskOverallocate   int64  `json:"disk_overallocate,omitempty" validate:"omitempty,gte=-1"`
	UploadSize         int64  `json:"upload_size,omitempty" validate:"omitempty,gt=0"`
	DaemonSFTP         int    `json:"daemon_sftp,omitempty" validate:"omitempty,gt=0,lt=65536"`
	DaemonListen       int    `json:"daemon_listen,omitempty" validate:"omitempty,gt=0,lt=65536"`
	DaemonBase         string `json:"daemon_base,omitempty"`
	Public             *bool  `json:"public,omitempty"`
	BehindProxy        *bool  `json:"behind_proxy,omitempty"`
	MaintenanceMode    *bool  `json:"maintenance_mode,omitempty"`
	Description        string `json:"description,omitempty"`
}

// CreateNodeAllocationRequest defines the request body for creating node allocations.
type CreateNodeAllocationRequest struct {
	IP    string   `json:"ip" validate:"required,ip"`
	Ports []string `json:"ports" validate:"required,min=1,dive,required"`
}

// ListNodes retrieves a paginated list of all nodes.
func (c *client) ListNodes(ctx context.Context, options pagination.ListOptions) ([]*models.Node, *pagination.Paginator[*models.Node], error) {
	return pagination.New[*models.Node](ctx, c.client, "application/nodes", options)
}

// GetNode retrieves details for a specific node by its ID.
func (c *client) GetNode(ctx context.Context, id int) (*models.Node, error) {
	path := fmt.Sprintf("application/nodes/%d", id)
	var response struct {
		Attributes models.Node `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// CreateNode creates a new node.
func (c *client) CreateNode(ctx context.Context, req CreateNodeRequest) (*models.Node, error) {
	var response struct {
		Attributes models.Node `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "application/nodes", req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// UpdateNode updates an existing node's configuration.
func (c *client) UpdateNode(ctx context.Context, id int, req UpdateNodeRequest) (*models.Node, error) {
	path := fmt.Sprintf("application/nodes/%d", id)
	var response struct {
		Attributes models.Node `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DeleteNode permanently deletes a node.
func (c *client) DeleteNode(ctx context.Context, id int) error {
	path := fmt.Sprintf("application/nodes/%d", id)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}

// GetNodeConfiguration retrieves the Wings configuration for a node.
func (c *client) GetNodeConfiguration(ctx context.Context, id int) (*models.NodeConfiguration, error) {
	path := fmt.Sprintf("application/nodes/%d/configuration", id)
	var config models.NodeConfiguration
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ListNodeAllocations retrieves all allocations for a specific node.
func (c *client) ListNodeAllocations(ctx context.Context, nodeID int, options pagination.ListOptions) ([]*models.Allocation, *pagination.Paginator[*models.Allocation], error) {
	path := fmt.Sprintf("application/nodes/%d/allocations", nodeID)
	return pagination.New[*models.Allocation](ctx, c.client, path, options)
}

// CreateNodeAllocations creates new allocations for a node.
func (c *client) CreateNodeAllocations(ctx context.Context, nodeID int, req CreateNodeAllocationRequest) error {
	path := fmt.Sprintf("application/nodes/%d/allocations", nodeID)
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// DeleteNodeAllocation deletes a specific allocation from a node.
func (c *client) DeleteNodeAllocation(ctx context.Context, nodeID, allocationID int) error {
	path := fmt.Sprintf("application/nodes/%d/allocations/%d", nodeID, allocationID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}

// GetDeployableNodes retrieves nodes that can accept a new server deployment.
func (c *client) GetDeployableNodes(ctx context.Context, memory, disk int64) ([]*models.Node, error) {
	path := fmt.Sprintf("application/nodes/deployable?memory=%d&disk=%d", memory, disk)
	var response struct {
		Data []struct {
			Attributes models.Node `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	nodes := make([]*models.Node, len(response.Data))
	for i, item := range response.Data {
		nodes[i] = &item.Attributes
	}
	return nodes, nil
}
