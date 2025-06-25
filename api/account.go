package api

import "time"

type Account struct {
	ID        int    `json:"id"`
	Admin     bool   `json:"admin"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Language  string `json:"language"`
}

type UpdateEmailOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordOptions struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"password"`
	PasswordConfirm string `json:"password_confirmation"`
}

type TwoFactorDetails struct {
	ImageURL string `json:"image_url_data"`
	Secret   string `json:"secret"`
}

type TwoFactorEnableOptions struct {
	Code string `json:"code"`
}

type TwoFactorDisableOptions struct {
	Password string `json:"password"`
}

// APIKey represents a client API key.
type APIKey struct {
	Identifier  string     `json:"identifier"`
	Description string     `json:"description"`
	AllowedIPs  []string   `json:"allowed_ips"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	// The secret token is only provided on creation.
	Token *string `json:"token,omitempty"`
}

// APIKeyCreateOptions defines the request body for creating a new API key.
type APIKeyCreateOptions struct {
	Description string   `json:"description"`
	AllowedIPs  []string `json:"allowed_ips"`
}

type APIKeyCreateResponse struct {
	Object     string `json:"object"`
	Attributes APIKey `json:"attributes"`
	Meta       struct {
		SecretToken string `json:"secret_token"`
	} `json:"meta"`
}
