package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// GetAccount retrieves information about the authenticated account.
func (c *client) GetAccount(ctx context.Context) (*models.User, error) {
	var response struct {
		Attributes models.User `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, "client/account", nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// GetTwoFactorQR generates a TOTP QR code for setting up two-factor authentication.
func (c *client) GetTwoFactorQR(ctx context.Context) (*models.TwoFactorData, error) {
	var response struct {
		Data models.TwoFactorData `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, "client/account/two-factor", nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}

// EnableTwoFactor enables two-factor authentication using a TOTP code.
func (c *client) EnableTwoFactor(ctx context.Context, code string) (*models.RecoveryTokens, error) {
	req := map[string]string{"code": code}
	var response struct {
		Attributes models.RecoveryTokens `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "client/account/two-factor", req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DisableTwoFactor disables two-factor authentication for the account.
func (c *client) DisableTwoFactor(ctx context.Context, password string) error {
	req := map[string]string{"password": password}
	_, err := c.client.Do(ctx, http.MethodPost, "client/account/two-factor/disable", req, nil)
	return err
}

// UpdateEmail updates the email address for the account.
func (c *client) UpdateEmail(ctx context.Context, email, password string) error {
	req := map[string]string{"email": email, "password": password}
	_, err := c.client.Do(ctx, http.MethodPut, "client/account/email", req, nil)
	return err
}

// UpdatePassword changes the account password.
func (c *client) UpdatePassword(ctx context.Context, currentPassword, newPassword, confirmPassword string) error {
	req := map[string]string{
		"current_password":      currentPassword,
		"password":              newPassword,
		"password_confirmation": confirmPassword,
	}
	_, err := c.client.Do(ctx, http.MethodPut, "client/account/password", req, nil)
	return err
}

// ListAPIKeys retrieves a list of all API keys for the account.
func (c *client) ListAPIKeys(ctx context.Context) ([]*models.APIKey, error) {
	var response struct {
		Data []struct {
			Attributes models.APIKey `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, "client/account/api-keys", nil, &response)
	if err != nil {
		return nil, err
	}

	keys := make([]*models.APIKey, len(response.Data))
	for i, item := range response.Data {
		keys[i] = &item.Attributes
	}
	return keys, nil
}

// CreateAPIKey generates a new API key.
func (c *client) CreateAPIKey(ctx context.Context, description string, allowedIPs []string) (*models.APIKey, error) {
	req := map[string]interface{}{
		"description": description,
		"allowed_ips": allowedIPs,
	}
	var response struct {
		Attributes models.APIKey `json:"attributes"`
		Meta       struct {
			SecretToken string `json:"secret_token"`
		} `json:"meta"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "client/account/api-keys", req, &response)
	if err != nil {
		return nil, err
	}
	response.Attributes.Meta.SecretToken = response.Meta.SecretToken
	return &response.Attributes, nil
}

// DeleteAPIKey deletes an API key by its identifier.
func (c *client) DeleteAPIKey(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("client/account/api-keys/%s", identifier)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
