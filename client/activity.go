package client

import (
	"context"
	"fmt"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// ListAccountActivity retrieves activity logs for the authenticated account.
func (c *client) ListAccountActivity(ctx context.Context, options pagination.ListOptions) ([]*models.ActivityLog, *pagination.Paginator[*models.ActivityLog], error) {
	return pagination.New[*models.ActivityLog](ctx, c.client, "client/account/activity", options)
}

// ListServerActivity retrieves activity logs for a specific server.
func (c *client) ListServerActivity(ctx context.Context, serverID string, options pagination.ListOptions) ([]*models.ActivityLog, *pagination.Paginator[*models.ActivityLog], error) {
	path := fmt.Sprintf("client/servers/%s/activity", serverID)
	return pagination.New[*models.ActivityLog](ctx, c.client, path, options)
}
