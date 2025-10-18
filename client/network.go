package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// ListAllocations retrieves all network allocations for a server.
func (c *client) ListAllocations(ctx context.Context, serverID string) ([]*models.Allocation, error) {
	path := fmt.Sprintf("client/servers/%s/network/allocations", serverID)
	var response struct {
		Data []struct {
			Attributes models.Allocation `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	allocations := make([]*models.Allocation, len(response.Data))
	for i, item := range response.Data {
		allocations[i] = &item.Attributes
	}
	return allocations, nil
}

// AssignAllocation assigns an available allocation to the server.
func (c *client) AssignAllocation(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("client/servers/%s/network/allocations", serverID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	return err
}

// SetPrimaryAllocation sets an allocation as the server's primary connection port.
func (c *client) SetPrimaryAllocation(ctx context.Context, serverID string, allocationID int) error {
	path := fmt.Sprintf("client/servers/%s/network/allocations/%d/primary", serverID, allocationID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	return err
}

// UpdateAllocationNotes adds or modifies notes for an allocation.
func (c *client) UpdateAllocationNotes(ctx context.Context, serverID string, allocationID int, notes string) error {
	path := fmt.Sprintf("client/servers/%s/network/allocations/%d", serverID, allocationID)
	req := map[string]string{"notes": notes}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// DeleteAllocation unassigns an allocation from the server.
func (c *client) DeleteAllocation(ctx context.Context, serverID string, allocationID int) error {
	path := fmt.Sprintf("client/servers/%s/network/allocations/%d", serverID, allocationID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
