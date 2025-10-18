// Package models contains the data structures for the Pterodactyl API.
package models

import (
	"encoding/json"
	"time"
)

// Meta contains metadata about a response, such as pagination.
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination contains pagination information for a list of resources.
type Pagination struct {
	Total       int      `json:"total"`
	Count       int      `json:"count"`
	PerPage     int      `json:"per_page"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
	Links       struct{} `json:"links"` // Links can be complex, often not needed
}

// User represents a user account on the panel.
type User struct {
	ID         int       `json:"id"`
	ExternalID *string   `json:"external_id,omitempty"`
	UUID       string    `json:"uuid"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Language   string    `json:"language"`
	RootAdmin  bool      `json:"root_admin"`
	TwoFactor  bool      `json:"2fa"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Servers    *[]Server `json:"servers,omitempty"` // Included relation
}

// Server represents a game server instance.
type Server struct {
	ID              int                `json:"id"`
	InternalID      int                `json:"internal_id,omitempty"` // From client API
	ExternalID      *string            `json:"external_id,omitempty"`
	UUID            string             `json:"uuid"`
	Identifier      string             `json:"identifier"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Status          *string            `json:"status"`
	Suspended       bool               `json:"is_suspended"`
	Installing      bool               `json:"is_installing"`
	Transferring    bool               `json:"is_transferring"`
	ServerOwner     bool               `json:"server_owner"`
	UserID          int                `json:"user"`
	NodeID          int                `json:"node"`
	AllocationID    int                `json:"allocation"`
	NestID          int                `json:"nest"`
	EggID           int                `json:"egg"`
	Invocation      string             `json:"invocation,omitempty"`
	DockerImage     string             `json:"docker_image"`
	EggFeatures     []string           `json:"egg_features,omitempty"`
	Limits          Limits             `json:"limits"`
	FeatureLimits   FeatureLimits      `json:"feature_limits"`
	Relationships   *Relationships     `json:"relationships,omitempty"` // For client API
	UserPermissions []string           `json:"user_permissions,omitempty"`
	SftpDetails     SftpDetails        `json:"sftp_details,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Node            string             `json:"node_name,omitempty"` // From client API list
	User            *User              `json:"-"`                   // Included relation
	Allocations     *[]Allocation      `json:"allocations,omitempty"`
	Variables       *[]StartupVariable `json:"variables,omitempty"`
}

// Relationships contains related data for a server.
type Relationships struct {
	Allocations PaginatedAllocations `json:"allocations"`
}

// PaginatedAllocations is a helper struct for nested allocation lists.
type PaginatedAllocations struct {
	Object string       `json:"object"`
	Data   []Allocation `json:"data"`
}

// Limits defines the resource limits for a server.
type Limits struct {
	Memory      int64   `json:"memory"`
	Swap        int64   `json:"swap"`
	Disk        int64   `json:"disk"`
	IO          int64   `json:"io"`
	CPU         int64   `json:"cpu"`
	Threads     *string `json:"threads,omitempty"`
	OOMDisabled *bool   `json:"oom_disabled,omitempty"`
}

// FeatureLimits defines the limits on features like databases and backups.
type FeatureLimits struct {
	Databases   int `json:"databases"`
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
}

// SftpDetails contains connection information for SFTP.
type SftpDetails struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// Allocation represents a network allocation (IP and port).
type Allocation struct {
	ID        int     `json:"id"`
	IP        string  `json:"ip"`
	IPAlias   *string `json:"ip_alias,omitempty"`
	Port      int     `json:"port"`
	Notes     *string `json:"notes,omitempty"`
	IsDefault bool    `json:"is_default"`
	ServerID  int     `json:"server,omitempty"` // from application API node allocation list
}

// APIKey represents a user's API key.
// The raw key is only provided on creation.
type APIKey struct {
	Identifier  string     `json:"identifier"`
	Description string     `json:"description"`
	AllowedIPs  []string   `json:"allowed_ips"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Meta        struct {
		SecretToken string `json:"secret_token"`
	} `json:"meta,omitempty"`
}

// SignedURL represents a pre-signed URL for file operations.
type SignedURL struct {
	URL string `json:"url"`
}

// FileObject represents a file or directory.
type FileObject struct {
	Name       string    `json:"name"`
	Mode       string    `json:"mode"`
	ModeBits   string    `json:"mode_bits"`
	Size       int64     `json:"size"`
	IsFile     bool      `json:"is_file"`
	IsSymlink  bool      `json:"is_symlink"`
	Mimetype   string    `json:"mimetype"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// Stats represents the real-time resource usage of a server.
type Stats struct {
	CurrentState string    `json:"current_state"`
	IsSuspended  bool      `json:"is_suspended"`
	Resources    Resources `json:"resources"`
}

// Resources contains the detailed resource usage statistics.
type Resources struct {
	MemoryBytes      int64   `json:"memory_bytes"`
	MemoryLimitBytes int64   `json:"memory_limit_bytes"`
	CPUAbsolute      float64 `json:"cpu_absolute"`
	DiskBytes        int64   `json:"disk_bytes"`
	NetworkRxBytes   int64   `json:"network_rx_bytes"`
	NetworkTxBytes   int64   `json:"network_tx_bytes"`
	Uptime           int64   `json:"uptime"`
}

// StartupVariable represents a server's startup environment variable.
type StartupVariable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	EnvVariable  string `json:"env_variable"`
	DefaultValue string `json:"default_value"`
	ServerValue  string `json:"server_value"`
	IsEditable   bool   `json:"is_editable"`
	Rules        string `json:"rules"`
}

// TwoFactorData holds the data for enabling 2FA.
type TwoFactorData struct {
	ImageURLData string `json:"image_url_data"`
	Secret       string `json:"secret"`
}

// RecoveryTokens holds the 2FA recovery tokens.
type RecoveryTokens struct {
	Tokens []string `json:"tokens"`
}

// Database represents a server database.
type Database struct {
	ID              string `json:"id"`
	Host            DatabaseHost
	Name            string `json:"name"`
	Username        string `json:"username"`
	ConnectionsFrom string `json:"connections_from"`
	MaxConnections  int    `json:"max_connections"`
	Password        *struct {
		Attributes DatabasePassword `json:"attributes"`
	} `json:"relationships,omitempty"`
}

// DatabaseHost contains connection info for the database host.
type DatabaseHost struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// DatabasePassword contains the password for a database.
// This is only returned on creation or password rotation.
type DatabasePassword struct {
	Password string `json:"password"`
}

// Backup represents a server backup.
type Backup struct {
	UUID         string     `json:"uuid"`
	Name         string     `json:"name"`
	IgnoredFiles []string   `json:"ignored_files"`
	SHA256Hash   *string    `json:"sha256_hash,omitempty"`
	Bytes        int64      `json:"bytes"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	IsSuccessful bool       `json:"is_successful"`
	IsLocked     bool       `json:"is_locked"`
}

// Subuser represents a user with access to a server.
type Subuser struct {
	UUID         string    `json:"uuid"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Image        string    `json:"image"`
	TwoFAEnabled bool      `json:"2fa_enabled"`
	CreatedAt    time.Time `json:"created_at"`
	Permissions  []string  `json:"permissions"`
}

// ActivityLog represents an entry in the activity log.
type ActivityLog struct {
	ID                    string          `json:"id"`
	Batch                 *string         `json:"batch"`
	Event                 string          `json:"event"`
	IsAPI                 bool            `json:"is_api"`
	IP                    string          `json:"ip"`
	Description           string          `json:"description"`
	Properties            json.RawMessage `json:"properties"`
	HasAdditionalMetadata bool            `json:"has_additional_metadata"`
	Timestamp             time.Time       `json:"timestamp"`
}

// StartupConfiguration represents the startup settings for a server.
type StartupConfiguration struct {
	Variables         []*StartupVariable
	StartupCommand    string
	RawStartupCommand string
	DockerImages      map[string]string
}

// SSHKey represents an SSH public key associated with an account.
type SSHKey struct {
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	PublicKey   string    `json:"public_key"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemPermissions represents the available permissions on the system.
type SystemPermissions struct {
	Permissions map[string]PermissionGroup `json:"permissions"`
}

// PermissionGroup is a group of related permissions.
type PermissionGroup struct {
	Description string            `json:"description"`
	Keys        map[string]string `json:"keys"`
}

// Schedule represents a scheduled task for a server.
type Schedule struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Cron           Cron            `json:"cron"`
	IsActive       bool            `json:"is_active"`
	IsProcessing   bool            `json:"is_processing"`
	OnlyWhenOnline bool            `json:"only_when_online"`
	LastRunAt      *time.Time      `json:"last_run_at"`
	NextRunAt      time.Time       `json:"next_run_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Tasks          *[]ScheduleTask `json:"tasks,omitempty"`
}

// Cron represents the cron timing for a schedule.
type Cron struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth string `json:"day_of_month"`
	Month      string `json:"month"`
	Hour       string `json:"hour"`
	Minute     string `json:"minute"`
}

// ScheduleTask represents a single task within a schedule.
type ScheduleTask struct {
	ID                int       `json:"id"`
	SequenceID        int       `json:"sequence_id"`
	Action            string    `json:"action"`
	Payload           string    `json:"payload"`
	TimeOffset        int       `json:"time_offset"`
	IsQueued          bool      `json:"is_queued"`
	ContinueOnFailure bool      `json:"continue_on_failure"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Node represents a Wings daemon node.
type Node struct {
	ID                 int       `json:"id"`
	UUID               string    `json:"uuid"`
	Public             bool      `json:"public"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	LocationID         int       `json:"location_id"`
	FQDN               string    `json:"fqdn"`
	Scheme             string    `json:"scheme"`
	BehindProxy        bool      `json:"behind_proxy"`
	MaintenanceMode    bool      `json:"maintenance_mode"`
	Memory             int64     `json:"memory"`
	MemoryOverallocate int64     `json:"memory_overallocate"`
	Disk               int64     `json:"disk"`
	DiskOverallocate   int64     `json:"disk_overallocate"`
	UploadSize         int64     `json:"upload_size"`
	DaemonListen       int       `json:"daemon_listen"`
	DaemonSFTP         int       `json:"daemon_sftp"`
	DaemonBase         string    `json:"daemon_base"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// NodeConfiguration represents the Wings configuration for a node.
type NodeConfiguration struct {
	Debug         bool             `json:"debug"`
	UUID          string           `json:"uuid"`
	TokenID       string           `json:"token_id"`
	Token         string           `json:"token"`
	API           NodeConfigAPI    `json:"api"`
	System        NodeConfigSystem `json:"system"`
	AllowedMounts []string         `json:"allowed_mounts"`
	Remote        string           `json:"remote"`
}

// NodeConfigAPI represents the API configuration section.
type NodeConfigAPI struct {
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	SSL         NodeConfigSSL `json:"ssl"`
	UploadLimit int           `json:"upload_limit"`
}

// NodeConfigSSL represents SSL configuration.
type NodeConfigSSL struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"cert"`
	Key     string `json:"key"`
}

// NodeConfigSystem represents system configuration.
type NodeConfigSystem struct {
	RootDirectory  string                   `json:"root_directory"`
	LogDirectory   string                   `json:"log_directory"`
	Data           string                   `json:"data"`
	SFTP           NodeConfigSFTP           `json:"sftp"`
	CrashDetection NodeConfigCrashDetection `json:"crash_detection"`
	Backups        NodeConfigBackups        `json:"backups"`
	Transfers      NodeConfigTransfers      `json:"transfers"`
}

// NodeConfigSFTP represents SFTP configuration.
type NodeConfigSFTP struct {
	BindPort int `json:"bind_port"`
}

// NodeConfigCrashDetection represents crash detection configuration.
type NodeConfigCrashDetection struct {
	Enabled bool `json:"enabled"`
	Timeout int  `json:"timeout"`
}

// NodeConfigBackups represents backup configuration.
type NodeConfigBackups struct {
	WriteLimit int `json:"write_limit"`
}

// NodeConfigTransfers represents transfer configuration.
type NodeConfigTransfers struct {
	DownloadLimit int `json:"download_limit"`
}

// Location represents a geographic location for organizing nodes.
type Location struct {
	ID        int       `json:"id"`
	Short     string    `json:"short"`
	Long      string    `json:"long"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Nest represents a server type category (e.g., Minecraft, Voice Servers).
type Nest struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	Author      string    `json:"author"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Egg represents a specific server configuration within a nest.
type Egg struct {
	ID           int               `json:"id"`
	UUID         string            `json:"uuid"`
	Name         string            `json:"name"`
	Nest         int               `json:"nest"`
	Author       string            `json:"author"`
	Description  string            `json:"description"`
	DockerImage  string            `json:"docker_image"`
	DockerImages map[string]string `json:"docker_images"`
	Config       EggConfig         `json:"config"`
	Startup      string            `json:"startup"`
	Script       EggScript         `json:"script"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Variables    *[]EggVariable    `json:"relationships,omitempty"`
}

// EggConfig represents the configuration section of an egg.
type EggConfig struct {
	Files   map[string]EggConfigFile `json:"files"`
	Startup EggConfigStartup         `json:"startup"`
	Stop    string                   `json:"stop"`
	Logs    map[string]interface{}   `json:"logs"`
	Extends *string                  `json:"extends,omitempty"`
}

// EggConfigFile represents a configuration file template.
type EggConfigFile struct {
	Parser string                 `json:"parser"`
	Find   map[string]interface{} `json:"find"`
}

// EggConfigStartup represents startup configuration.
type EggConfigStartup struct {
	Done            string   `json:"done"`
	UserInteraction []string `json:"user_interaction"`
}

// EggScript represents the installation script configuration.
type EggScript struct {
	Privileged bool    `json:"privileged"`
	Install    string  `json:"install"`
	Entry      string  `json:"entry"`
	Container  string  `json:"container"`
	Extends    *string `json:"extends,omitempty"`
}

// EggVariable represents a configuration variable for an egg.
type EggVariable struct {
	ID           int       `json:"id"`
	EggID        int       `json:"egg_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	EnvVariable  string    `json:"env_variable"`
	DefaultValue string    `json:"default_value"`
	UserViewable bool      `json:"user_viewable"`
	UserEditable bool      `json:"user_editable"`
	Rules        string    `json:"rules"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ApplicationDatabaseHost represents a database host in the application API.
type ApplicationDatabaseHost struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	NodeID    int       `json:"node"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ApplicationDatabase represents a database in the application API.
type ApplicationDatabase struct {
	ID             int       `json:"id"`
	ServerID       int       `json:"server"`
	HostID         int       `json:"host"`
	Database       string    `json:"database"`
	Username       string    `json:"username"`
	Remote         string    `json:"remote"`
	MaxConnections int       `json:"max_connections"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
