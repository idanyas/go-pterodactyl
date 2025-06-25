package api

type Resources struct {
	Object     string             `json:"object"`
	Attributes ResourceAttributes `json:"attributes"`
}

// ResourceAttributes contains the detailed usage statistics.
type ResourceAttributes struct {
	CurrentState string `json:"current_state"` // e.g., "running", "offline"
	IsSuspended  bool   `json:"is_suspended"`
	Resources    struct {
		MemoryBytes    int     `json:"memory_bytes"`
		CPUAbsolute    float64 `json:"cpu_absolute"`
		DiskBytes      int     `json:"disk_bytes"`
		NetworkRxBytes int     `json:"network_rx_bytes"`
		NetworkTxBytes int     `json:"network_tx_bytes"`
	} `json:"resources"`
}
