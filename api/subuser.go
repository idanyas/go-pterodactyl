package api

import "time"

type Subuser struct {
	UUID             string    `json:"uuid"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Image            string    `json:"image"`
	TwoFactorEnabled bool      `json:"2fa_enabled"`
	CreatedAt        time.Time `json:"created_at"`
	// A list of permission strings that this subuser has.
	Permissions []string `json:"permissions"`
}

// SubuserCreateOptions defines the request body for creating a new subuser.
type SubuserCreateOptions struct {
	// Email address of the user to add.
	Email string `json:"email"`
	// A list of permission strings to grant to the user.
	Permissions []string `json:"permissions"`
}

// SubuserUpdateOptions defines the request body for updating a subuser's permissions.
type SubuserUpdateOptions struct {
	// A list of permission strings to grant to the user.
	// This will overwrite all existing permissions.
	Permissions []string `json:"permissions"`
}

// subuserResponse is a helper for unmarshaling single subuser responses.
type SubuserResponse struct {
	Object     string   `json:"object"`
	Attributes *Subuser `json:"attributes"`
}
