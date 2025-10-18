package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// ListNests retrieves a paginated list of all nests.
func (c *client) ListNests(ctx context.Context, options pagination.ListOptions) ([]*models.Nest, *pagination.Paginator[*models.Nest], error) {
	return pagination.New[*models.Nest](ctx, c.client, "application/nests", options)
}

// GetNest retrieves details for a specific nest by its ID.
func (c *client) GetNest(ctx context.Context, id int) (*models.Nest, error) {
	path := fmt.Sprintf("application/nests/%d", id)
	var response struct {
		Attributes models.Nest `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// ListNestEggs retrieves all eggs within a specific nest.
func (c *client) ListNestEggs(ctx context.Context, nestID int, options pagination.ListOptions) ([]*models.Egg, *pagination.Paginator[*models.Egg], error) {
	path := fmt.Sprintf("application/nests/%d/eggs", nestID)
	return pagination.New[*models.Egg](ctx, c.client, path, options)
}

// GetEgg retrieves details for a specific egg.
func (c *client) GetEgg(ctx context.Context, nestID, eggID int) (*models.Egg, error) {
	path := fmt.Sprintf("application/nests/%d/eggs/%d", nestID, eggID)
	var response struct {
		Attributes models.Egg `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}
