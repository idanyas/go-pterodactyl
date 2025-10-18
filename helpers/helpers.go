// Package helpers provides convenience methods for common Pterodactyl operations.
package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/idanyas/go-pterodactyl/client"
	"github.com/idanyas/go-pterodactyl/models"
)

// StateWaiter provides methods to wait for server state changes.
type StateWaiter struct {
	client client.ClientClient
}

// NewStateWaiter creates a new StateWaiter.
func NewStateWaiter(c client.ClientClient) *StateWaiter {
	return &StateWaiter{client: c}
}

// WaitForState polls the server until it reaches the desired state or the context times out.
// Returns an error if the context is cancelled or if polling fails.
func (w *StateWaiter) WaitForState(ctx context.Context, serverID string, desiredState string, pollInterval time.Duration) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resources, err := w.client.GetServerResources(ctx, serverID)
			if err != nil {
				return fmt.Errorf("failed to get server resources: %w", err)
			}

			if resources.CurrentState == desiredState {
				return nil
			}
		}
	}
}

// FileDownloader provides methods for downloading server files.
type FileDownloader struct {
	client client.ClientClient
}

// NewFileDownloader creates a new FileDownloader.
func NewFileDownloader(c client.ClientClient) *FileDownloader {
	return &FileDownloader{client: c}
}

// DownloadToWriter downloads a file from a server and writes it to the provided writer.
func (d *FileDownloader) DownloadToWriter(ctx context.Context, serverID, filePath string, w io.Writer) error {
	// Get the download URL
	signedURL, err := d.client.GetDownloadURL(ctx, serverID, filePath)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, signedURL.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	httpClient := &http.Client{Timeout: 5 * time.Minute}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Copy to writer
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// BackupManager provides high-level backup operations.
type BackupManager struct {
	client client.ClientClient
}

// NewBackupManager creates a new BackupManager.
func NewBackupManager(c client.ClientClient) *BackupManager {
	return &BackupManager{client: c}
}

// CreateAndWait creates a backup and waits for it to complete.
func (m *BackupManager) CreateAndWait(ctx context.Context, serverID string, req client.CreateBackupRequest, pollInterval time.Duration) (*models.Backup, error) {
	// Create the backup
	backup, err := m.client.CreateBackup(ctx, serverID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Poll until complete
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			current, err := m.client.GetBackup(ctx, serverID, backup.UUID)
			if err != nil {
				return nil, fmt.Errorf("failed to get backup status: %w", err)
			}

			if current.CompletedAt != nil {
				return current, nil
			}
		}
	}
}
