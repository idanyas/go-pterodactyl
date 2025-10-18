// Package client provides the Client API implementation for the Pterodactyl Panel.
package client

import (
	"context"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
	"github.com/idanyas/go-pterodactyl/websocket"
)

// PaginatorClient defines the interface required by the Paginator.
type PaginatorClient interface {
	Do(ctx context.Context, method, path string, body, v interface{}) (*http.Response, error)
}

// ClientClient is the client for the Client API.
type ClientClient interface {
	// Account Management
	GetAccount(ctx context.Context) (*models.User, error)
	GetTwoFactorQR(ctx context.Context) (*models.TwoFactorData, error)
	EnableTwoFactor(ctx context.Context, code string) (*models.RecoveryTokens, error)
	DisableTwoFactor(ctx context.Context, password string) error
	UpdateEmail(ctx context.Context, email, password string) error
	UpdatePassword(ctx context.Context, currentPassword, newPassword, confirmPassword string) error
	ListAPIKeys(ctx context.Context) ([]*models.APIKey, error)
	CreateAPIKey(ctx context.Context, description string, allowedIPs []string) (*models.APIKey, error)
	DeleteAPIKey(ctx context.Context, identifier string) error

	// SSH Key Management
	ListSSHKeys(ctx context.Context) ([]*models.SSHKey, error)
	AddSSHKey(ctx context.Context, name, publicKey string) (*models.SSHKey, error)
	RemoveSSHKey(ctx context.Context, fingerprint string) error

	// Activity Logs
	ListAccountActivity(ctx context.Context, options pagination.ListOptions) ([]*models.ActivityLog, *pagination.Paginator[*models.ActivityLog], error)
	ListServerActivity(ctx context.Context, serverID string, options pagination.ListOptions) ([]*models.ActivityLog, *pagination.Paginator[*models.ActivityLog], error)

	// Permissions
	GetSystemPermissions(ctx context.Context) (*models.SystemPermissions, error)

	// Server Management
	ListServers(ctx context.Context, options pagination.ListOptions) ([]*models.Server, *pagination.Paginator[*models.Server], error)
	GetServer(ctx context.Context, serverID string) (*models.Server, error)
	GetServerResources(ctx context.Context, serverID string) (*models.Stats, error)
	SendPowerAction(ctx context.Context, serverID, signal string) error
	SendCommand(ctx context.Context, serverID, command string) error

	// WebSocket
	ConnectWebSocket(ctx context.Context, serverID string) (*websocket.Conn, error)

	// File Management
	ListFiles(ctx context.Context, serverID, directory string) ([]*models.FileObject, error)
	GetFileContents(ctx context.Context, serverID, filePath string) (string, error)
	WriteFile(ctx context.Context, serverID, filePath, content string) error
	CreateDirectory(ctx context.Context, serverID, root, name string) error
	DeleteFiles(ctx context.Context, serverID, root string, files []string) error
	RenameFile(ctx context.Context, serverID, root, from, to string) error
	CopyFile(ctx context.Context, serverID, location string) error
	GetDownloadURL(ctx context.Context, serverID, filePath string) (*models.SignedURL, error)
	GetUploadURL(ctx context.Context, serverID, directory string) (*models.SignedURL, error)
	CompressFiles(ctx context.Context, serverID, root string, files []string) (*models.FileObject, error)
	DecompressFile(ctx context.Context, serverID, root, file string) error
	ChmodFiles(ctx context.Context, serverID, root string, files []ChmodFileRequest) error
	PullFile(ctx context.Context, serverID, url, directory, filename string) error

	// Database Management
	ListDatabases(ctx context.Context, serverID string) ([]*models.Database, error)
	CreateDatabase(ctx context.Context, serverID, database, remote string) (*models.Database, error)
	RotateDatabasePassword(ctx context.Context, serverID, databaseID string) (*models.Database, error)
	DeleteDatabase(ctx context.Context, serverID, databaseID string) error

	// Backup Management
	ListBackups(ctx context.Context, serverID string) ([]*models.Backup, error)
	GetBackup(ctx context.Context, serverID, backupUUID string) (*models.Backup, error)
	CreateBackup(ctx context.Context, serverID string, req CreateBackupRequest) (*models.Backup, error)
	GetBackupDownloadURL(ctx context.Context, serverID, backupUUID string) (*models.SignedURL, error)
	RestoreBackup(ctx context.Context, serverID, backupUUID string, truncate bool) error
	ToggleBackupLock(ctx context.Context, serverID, backupUUID string) error
	DeleteBackup(ctx context.Context, serverID, backupUUID string) error

	// Startup Management
	GetStartupConfig(ctx context.Context, serverID string) (*models.StartupConfiguration, error)
	UpdateStartupVariable(ctx context.Context, serverID, key, value string) (*models.StartupVariable, error)

	// Settings
	RenameServer(ctx context.Context, serverID, name, description string) error
	ReinstallServer(ctx context.Context, serverID string) error
	UpdateDockerImage(ctx context.Context, serverID, dockerImage string) error

	// Network Management
	ListAllocations(ctx context.Context, serverID string) ([]*models.Allocation, error)
	AssignAllocation(ctx context.Context, serverID string) error
	SetPrimaryAllocation(ctx context.Context, serverID string, allocationID int) error
	UpdateAllocationNotes(ctx context.Context, serverID string, allocationID int, notes string) error
	DeleteAllocation(ctx context.Context, serverID string, allocationID int) error

	// Subuser Management
	ListSubusers(ctx context.Context, serverID string) ([]*models.Subuser, error)
	GetSubuser(ctx context.Context, serverID, userUUID string) (*models.Subuser, error)
	CreateSubuser(ctx context.Context, serverID, email string, permissions []string) (*models.Subuser, error)
	UpdateSubuser(ctx context.Context, serverID, userUUID string, permissions []string) (*models.Subuser, error)
	DeleteSubuser(ctx context.Context, serverID, userUUID string) error

	// Schedule Management
	ListSchedules(ctx context.Context, serverID string) ([]*models.Schedule, error)
	GetSchedule(ctx context.Context, serverID string, scheduleID int) (*models.Schedule, error)
	CreateSchedule(ctx context.Context, serverID string, req CreateScheduleRequest) (*models.Schedule, error)
	UpdateSchedule(ctx context.Context, serverID string, scheduleID int, req UpdateScheduleRequest) (*models.Schedule, error)
	DeleteSchedule(ctx context.Context, serverID string, scheduleID int) error
	ExecuteSchedule(ctx context.Context, serverID string, scheduleID int) error

	// Schedule Task Management
	CreateScheduleTask(ctx context.Context, serverID string, scheduleID int, req CreateScheduleTaskRequest) (*models.ScheduleTask, error)
	UpdateScheduleTask(ctx context.Context, serverID string, scheduleID, taskID int, req UpdateScheduleTaskRequest) (*models.ScheduleTask, error)
	DeleteScheduleTask(ctx context.Context, serverID string, scheduleID, taskID int) error
}

type client struct {
	client PaginatorClient
}

// New creates a new Client API client.
func New(c PaginatorClient) ClientClient {
	return &client{client: c}
}
