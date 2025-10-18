package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// CreateLocationRequest defines the request body for creating a new location.
type CreateLocationRequest struct {
	Short string `json:"short" validate:"required,min=1"`
	Long  string `json:"long,omitempty"`
}

// UpdateLocationRequest defines the request body for updating a location.
type UpdateLocationRequest struct {
	Short string `json:"short,omitempty" validate:"omitempty,min=1"`
	Long  string `json:"long,omitempty"`
}

// ListLocations retrieves a paginated list of all locations.
func (c *client) ListLocations(ctx context.Context, options pagination.ListOptions) ([]*models.Location, *pagination.Paginator[*models.Location], error) {
	return pagination.New[*models.Location](ctx, c.client, "application/locations", options)
}

// GetLocation retrieves details for a specific location by its ID.
func (c *client) GetLocation(ctx context.Context, id int) (*models.Location, error) {
	if id <= 0 {
		return nil, fmt.Errorf("location ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/locations/%d", id)
	var response struct {
		Attributes models.Location `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get location %d: %w", id, err)
	}
	return &response.Attributes, nil
}

// CreateLocation creates a new location.
func (c *client) CreateLocation(ctx context.Context, req CreateLocationRequest) (*models.Location, error) {
	var response struct {
		Attributes models.Location `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "application/locations", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}
	return &response.Attributes, nil
}

// UpdateLocation updates an existing location's details.
func (c *client) UpdateLocation(ctx context.Context, id int, req UpdateLocationRequest) (*models.Location, error) {
	if id <= 0 {
		return nil, fmt.Errorf("location ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/locations/%d", id)
	var response struct {
		Attributes models.Location `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update location %d: %w", id, err)
	}
	return &response.Attributes, nil
}

// DeleteLocation permanently deletes a location.
func (c *client) DeleteLocation(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("location ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/locations/%d", id)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete location %d: %w", id, err)
	}
	return nil
}
