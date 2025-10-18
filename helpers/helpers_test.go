package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/idanyas/go-pterodactyl/client"
	"github.com/idanyas/go-pterodactyl/models"
	"github.com/idanyas/go-pterodactyl/pagination"
	"github.com/idanyas/go-pterodactyl/websocket"
)

// mockClientForHelpers is a mock implementation of client.ClientClient for testing helpers.
type mockClientForHelpers struct{}

// Implemented methods for tests
func (m *mockClientForHelpers) GetServerResources(ctx context.Context, serverID string) (*models.Stats, error) {
	return &models.Stats{
		CurrentState: "running",
	}, nil
}

func (m *mockClientForHelpers) GetDownloadURL(ctx context.Context, serverID, filePath string) (*models.SignedURL, error) {
	return &models.SignedURL{
		URL: "https://example.com/download",
	}, nil
}

func (m *mockClientForHelpers) CreateBackup(ctx context.Context, serverID string, req client.CreateBackupRequest) (*models.Backup, error) {
	now := time.Now()
	return &models.Backup{
		UUID:        "test-uuid",
		Name:        "test",
		CompletedAt: &now,
	}, nil
}

func (m *mockClientForHelpers) GetBackup(ctx context.Context, serverID, backupUUID string) (*models.Backup, error) {
	now := time.Now()
	return &models.Backup{
		UUID:        backupUUID,
		Name:        "test",
		CompletedAt: &now,
	}, nil
}

// Stub implementations to satisfy the interface
func (m *mockClientForHelpers) GetAccount(ctx context.Context) (*models.User, error) { return nil, nil }
func (m *mockClientForHelpers) GetTwoFactorQR(ctx context.Context) (*models.TwoFactorData, error) {
	return nil, nil
}
func (m *mockClientForHelpers) EnableTwoFactor(ctx context.Context, code string) (*models.RecoveryTokens, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DisableTwoFactor(ctx context.Context, password string) error {
	return nil
}
func (m *mockClientForHelpers) UpdateEmail(ctx context.Context, email, password string) error {
	return nil
}
func (m *mockClientForHelpers) UpdatePassword(ctx context.Context, currentPassword, newPassword, confirmPassword string) error {
	return nil
}
func (m *mockClientForHelpers) ListAPIKeys(ctx context.Context) ([]*models.APIKey, error) {
	return nil, nil
}
func (m *mockClientForHelpers) CreateAPIKey(ctx context.Context, description string, allowedIPs []string) (*models.APIKey, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DeleteAPIKey(ctx context.Context, identifier string) error { return nil }
func (m *mockClientForHelpers) ListSSHKeys(ctx context.Context) ([]*models.SSHKey, error) {
	return nil, nil
}
func (m *mockClientForHelpers) AddSSHKey(ctx context.Context, name, publicKey string) (*models.SSHKey, error) {
	return nil, nil
}
func (m *mockClientForHelpers) RemoveSSHKey(ctx context.Context, fingerprint string) error {
	return nil
}
func (m *mockClientForHelpers) ListAccountActivity(ctx context.Context, options pagination.ListOptions) ([]*models.ActivityLog, *pagination.Paginator[*models.ActivityLog], error) {
	return nil, nil, nil
}
func (m *mockClientForHelpers) ListServerActivity(ctx context.Context, serverID string, options pagination.ListOptions) ([]*models.ActivityLog, *pagination.Paginator[*models.ActivityLog], error) {
	return nil, nil, nil
}
func (m *mockClientForHelpers) GetSystemPermissions(ctx context.Context) (*models.SystemPermissions, error) {
	return nil, nil
}
func (m *mockClientForHelpers) ListServers(ctx context.Context, options pagination.ListOptions) ([]*models.Server, *pagination.Paginator[*models.Server], error) {
	return nil, nil, nil
}
func (m *mockClientForHelpers) GetServer(ctx context.Context, serverID string) (*models.Server, error) {
	return nil, nil
}
func (m *mockClientForHelpers) SendPowerAction(ctx context.Context, serverID, signal string) error {
	return nil
}
func (m *mockClientForHelpers) SendCommand(ctx context.Context, serverID, command string) error {
	return nil
}
func (m *mockClientForHelpers) ConnectWebSocket(ctx context.Context, serverID string) (*websocket.Conn, error) {
	return nil, nil
}
func (m *mockClientForHelpers) ListFiles(ctx context.Context, serverID, directory string) ([]*models.FileObject, error) {
	return nil, nil
}
func (m *mockClientForHelpers) GetFileContents(ctx context.Context, serverID, filePath string) (string, error) {
	return "", nil
}
func (m *mockClientForHelpers) WriteFile(ctx context.Context, serverID, filePath, content string) error {
	return nil
}
func (m *mockClientForHelpers) CreateDirectory(ctx context.Context, serverID, root, name string) error {
	return nil
}
func (m *mockClientForHelpers) DeleteFiles(ctx context.Context, serverID, root string, files []string) error {
	return nil
}
func (m *mockClientForHelpers) RenameFile(ctx context.Context, serverID, root, from, to string) error {
	return nil
}
func (m *mockClientForHelpers) CopyFile(ctx context.Context, serverID, location string) error {
	return nil
}
func (m *mockClientForHelpers) GetUploadURL(ctx context.Context, serverID, directory string) (*models.SignedURL, error) {
	return nil, nil
}
func (m *mockClientForHelpers) CompressFiles(ctx context.Context, serverID, root string, files []string) (*models.FileObject, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DecompressFile(ctx context.Context, serverID, root, file string) error {
	return nil
}
func (m *mockClientForHelpers) ChmodFiles(ctx context.Context, serverID, root string, files []client.ChmodFileRequest) error {
	return nil
}
func (m *mockClientForHelpers) PullFile(ctx context.Context, serverID, url, directory, filename string) error {
	return nil
}
func (m *mockClientForHelpers) ListDatabases(ctx context.Context, serverID string) ([]*models.Database, error) {
	return nil, nil
}
func (m *mockClientForHelpers) CreateDatabase(ctx context.Context, serverID, database, remote string) (*models.Database, error) {
	return nil, nil
}
func (m *mockClientForHelpers) RotateDatabasePassword(ctx context.Context, serverID, databaseID string) (*models.Database, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DeleteDatabase(ctx context.Context, serverID, databaseID string) error {
	return nil
}
func (m *mockClientForHelpers) ListBackups(ctx context.Context, serverID string) ([]*models.Backup, error) {
	return nil, nil
}
func (m *mockClientForHelpers) GetBackupDownloadURL(ctx context.Context, serverID, backupUUID string) (*models.SignedURL, error) {
	return nil, nil
}
func (m *mockClientForHelpers) RestoreBackup(ctx context.Context, serverID, backupUUID string, truncate bool) error {
	return nil
}
func (m *mockClientForHelpers) ToggleBackupLock(ctx context.Context, serverID, backupUUID string) error {
	return nil
}
func (m *mockClientForHelpers) DeleteBackup(ctx context.Context, serverID, backupUUID string) error {
	return nil
}
func (m *mockClientForHelpers) GetStartupConfig(ctx context.Context, serverID string) (*models.StartupConfiguration, error) {
	return nil, nil
}
func (m *mockClientForHelpers) UpdateStartupVariable(ctx context.Context, serverID, key, value string) (*models.StartupVariable, error) {
	return nil, nil
}
func (m *mockClientForHelpers) RenameServer(ctx context.Context, serverID, name, description string) error {
	return nil
}
func (m *mockClientForHelpers) ReinstallServer(ctx context.Context, serverID string) error {
	return nil
}
func (m *mockClientForHelpers) UpdateDockerImage(ctx context.Context, serverID, dockerImage string) error {
	return nil
}
func (m *mockClientForHelpers) ListAllocations(ctx context.Context, serverID string) ([]*models.Allocation, error) {
	return nil, nil
}
func (m *mockClientForHelpers) AssignAllocation(ctx context.Context, serverID string) error {
	return nil
}
func (m *mockClientForHelpers) SetPrimaryAllocation(ctx context.Context, serverID string, allocationID int) error {
	return nil
}
func (m *mockClientForHelpers) UpdateAllocationNotes(ctx context.Context, serverID string, allocationID int, notes string) error {
	return nil
}
func (m *mockClientForHelpers) DeleteAllocation(ctx context.Context, serverID string, allocationID int) error {
	return nil
}
func (m *mockClientForHelpers) ListSubusers(ctx context.Context, serverID string) ([]*models.Subuser, error) {
	return nil, nil
}
func (m *mockClientForHelpers) GetSubuser(ctx context.Context, serverID, userUUID string) (*models.Subuser, error) {
	return nil, nil
}
func (m *mockClientForHelpers) CreateSubuser(ctx context.Context, serverID, email string, permissions []string) (*models.Subuser, error) {
	return nil, nil
}
func (m *mockClientForHelpers) UpdateSubuser(ctx context.Context, serverID, userUUID string, permissions []string) (*models.Subuser, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DeleteSubuser(ctx context.Context, serverID, userUUID string) error {
	return nil
}
func (m *mockClientForHelpers) ListSchedules(ctx context.Context, serverID string) ([]*models.Schedule, error) {
	return nil, nil
}
func (m *mockClientForHelpers) GetSchedule(ctx context.Context, serverID string, scheduleID int) (*models.Schedule, error) {
	return nil, nil
}
func (m *mockClientForHelpers) CreateSchedule(ctx context.Context, serverID string, req client.CreateScheduleRequest) (*models.Schedule, error) {
	return nil, nil
}
func (m *mockClientForHelpers) UpdateSchedule(ctx context.Context, serverID string, scheduleID int, req client.UpdateScheduleRequest) (*models.Schedule, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DeleteSchedule(ctx context.Context, serverID string, scheduleID int) error {
	return nil
}
func (m *mockClientForHelpers) ExecuteSchedule(ctx context.Context, serverID string, scheduleID int) error {
	return nil
}
func (m *mockClientForHelpers) CreateScheduleTask(ctx context.Context, serverID string, scheduleID int, req client.CreateScheduleTaskRequest) (*models.ScheduleTask, error) {
	return nil, nil
}
func (m *mockClientForHelpers) UpdateScheduleTask(ctx context.Context, serverID string, scheduleID, taskID int, req client.UpdateScheduleTaskRequest) (*models.ScheduleTask, error) {
	return nil, nil
}
func (m *mockClientForHelpers) DeleteScheduleTask(ctx context.Context, serverID string, scheduleID, taskID int) error {
	return nil
}
func (m *mockClientForHelpers) ConnectWebSocketWithReconnect(ctx context.Context, serverID string, reconnectOpts *websocket.ReconnectOptions) (*websocket.Conn, error) {
	return nil, nil
}

func TestStateWaiter_WaitForState(t *testing.T) {
	mock := &mockClientForHelpers{}
	// This would need a proper mock that changes state
	// For now, just test that it doesn't panic
	waiter := NewStateWaiter(mock)
	if waiter == nil {
		t.Error("NewStateWaiter returned nil")
	}
}

func TestFileDownloader(t *testing.T) {
	mock := &mockClientForHelpers{}
	downloader := NewFileDownloader(mock)
	if downloader == nil {
		t.Error("NewFileDownloader returned nil")
	}
}

func TestBackupManager(t *testing.T) {
	mock := &mockClientForHelpers{}
	manager := NewBackupManager(mock)
	if manager == nil {
		t.Error("NewBackupManager returned nil")
	}
}
