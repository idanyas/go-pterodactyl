package api

type PermissionDescriptor struct {
	Description string            `json:"description"`
	Keys        map[string]string `json:"keys"`
}

type Permission struct {
	Object     string                          `json:"object"`
	Attributes map[string]PermissionDescriptor `json:"attributes"`
}
