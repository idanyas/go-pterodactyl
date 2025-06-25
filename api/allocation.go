package api

type Allocation struct {
	ID       int     `json:"id"`
	IP       string  `json:"ip"`
	Alias    *string `json:"alias"`
	Port     int     `json:"port"`
	Notes    *string `json:"notes"`
	Assigned bool    `json:"assigned"`
}

type AllocationNoteOptions struct {
	Notes *string `json:"notes"`
}

type AllocationResponse struct {
	Object     string      `json:"object"`
	Attributes *Allocation `json:"attributes"`
}

type AllocationCreateOptions struct {
	IP    string   `json:"ip"`
	Ports []string `json:"ports"`
}
