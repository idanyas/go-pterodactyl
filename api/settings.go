package api

type RenameOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
