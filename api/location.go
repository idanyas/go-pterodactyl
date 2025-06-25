package api

import "time"

type Location struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"short"`
	Description string    `json:"long"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LocationCreateOptions struct {
	ShortCode   string `json:"short"`
	Description string `json:"long"`
}

type LocationUpdateOptions struct {
	ShortCode   string `json:"short,omitempty"`
	Description string `json:"long,omitempty"`
}
