package api

import "time"

type Database struct {
	ID             int       `json:"id"`
	ServerID       int       `json:"server"`
	HostID         int       `json:"host"`
	DatabaseName   string    `json:"database"`
	Username       string    `json:"username"`
	Remote         string    `json:"remote"`
	MaxConnections *int      `json:"max_connections"` // Use pointer for nullable number
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DatabaseCreateOptions struct {
	DatabaseName string `json:"database"`
	Remote       string `json:"remote"`
	// The Host field is optional and will be auto-assigned if omitted.
	Host *int `json:"host,omitempty"`
}

type ClientDatabase struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ConnectionsFrom string `json:"connections_from"`
	Host            struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"host"`
	// Password is not part of the standard API response for a database.
	// It is only returned on creation or password rotation, and our SDK
	// populates it manually. The json:"-" tag prevents it from being
	// marshaled or unmarshaled in normal operations.
	Password string `json:"-"`
}

// ClientDatabaseCreateOptions defines the request body for creating a new database.
type ClientDatabaseCreateOptions struct {
	DatabaseName string `json:"database"`
	Remote       string `json:"remote"`
}

// ClientDatabaseCreateResponse represents the unique response from creating or rotating
// a database password, which includes the secret in a relationships object.
type ClientDatabaseCreateResponse struct {
	Object        string          `json:"object"`
	Attributes    *ClientDatabase `json:"attributes"`
	Relationships *struct {
		Password *struct {
			Object     string `json:"object"`
			Attributes *struct {
				Password string `json:"password"`
			} `json:"attributes"`
		} `json:"password"`
	} `json:"relationships"`
}
