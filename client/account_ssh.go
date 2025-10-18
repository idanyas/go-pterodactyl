package client

import (
	"context"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// ListSSHKeys retrieves a list of SSH public keys associated with the account.
func (c *client) ListSSHKeys(ctx context.Context) ([]*models.SSHKey, error) {
	var response struct {
		Data []struct {
			Attributes models.SSHKey `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, "client/account/ssh-keys", nil, &response)
	if err != nil {
		return nil, err
	}

	keys := make([]*models.SSHKey, len(response.Data))
	for i, item := range response.Data {
		keys[i] = &item.Attributes
	}
	return keys, nil
}

// AddSSHKey adds a new SSH public key to the account.
func (c *client) AddSSHKey(ctx context.Context, name, publicKey string) (*models.SSHKey, error) {
	req := map[string]string{
		"name":       name,
		"public_key": publicKey,
	}
	var response struct {
		Attributes models.SSHKey `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, "client/account/ssh-keys", req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// RemoveSSHKey removes an SSH key from the account using its fingerprint.
func (c *client) RemoveSSHKey(ctx context.Context, fingerprint string) error {
	req := map[string]string{"fingerprint": fingerprint}
	_, err := c.client.Do(ctx, http.MethodPost, "client/account/ssh-keys/remove", req, nil)
	return err
}
