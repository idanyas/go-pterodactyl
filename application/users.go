package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
)

// CreateUserRequest defines the request body for creating a new user.
type CreateUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required,min=1,max=255"`
	FirstName  string `json:"first_name" validate:"required,min=1"`
	LastName   string `json:"last_name" validate:"required,min=1"`
	Password   string `json:"password,omitempty" validate:"omitempty,min=8"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
}

// UpdateUserRequest defines the request body for updating a user.
// All fields are optional.
type UpdateUserRequest struct {
	Email      string `json:"email,omitempty" validate:"omitempty,email"`
	Username   string `json:"username,omitempty" validate:"omitempty,min=1,max=255"`
	FirstName  string `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName   string `json:"last_name,omitempty" validate:"omitempty,min=1"`
	Password   string `json:"password,omitempty" validate:"omitempty,min=8"`
	RootAdmin  *bool  `json:"root_admin,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
}

// ListUsers retrieves a paginated list of all users.
func (c *client) ListUsers(ctx context.Context, options pagination.ListOptions) ([]*models.User, *pagination.Paginator[*models.User], error) {
	return pagination.New[*models.User](ctx, c.client, "application/users", options)
}

// GetUser retrieves details for a specific user by their ID.
func (c *client) GetUser(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("user ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/users/%d", id)
	var response struct {
		Attributes models.User `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %d: %w", id, err)
	}
	return &response.Attributes, nil
}

// GetUserExternal retrieves details for a specific user by their external ID.
func (c *client) GetUserExternal(ctx context.Context, externalID string) (*models.User, error) {
	if externalID == "" {
		return nil, fmt.Errorf("external ID cannot be empty")
	}

	path := fmt.Sprintf("application/users/external/%s", externalID)
	var response struct {
		Attributes models.User `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with external ID %s: %w", externalID, err)
	}
	return &response.Attributes, nil
}

// CreateUser creates a new user account.
func (c *client) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	var response struct {
		Attributes models.User `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "application/users", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &response.Attributes, nil
}

// UpdateUser updates an existing user's details.
func (c *client) UpdateUser(ctx context.Context, id int, req UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("user ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/users/%d", id)
	var response struct {
		Attributes models.User `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPatch, path, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update user %d: %w", id, err)
	}
	return &response.Attributes, nil
}

// DeleteUser permanently deletes a user account.
func (c *client) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("user ID must be positive, got %d", id)
	}

	path := fmt.Sprintf("application/users/%d", id)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete user %d: %w", id, err)
	}
	return nil
}
