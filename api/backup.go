package api

import "time"

type Backup struct {
	UUID         string     `json:"uuid"`
	IsSuccessful bool       `json:"is_successful"`
	IsLocked     bool       `json:"is_locked"`
	Name         string     `json:"name"`
	IgnoredFiles []string   `json:"ignored_files"`
	Checksum     *string    `json:"checksum"` // Can be null
	Bytes        int64      `json:"bytes"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"` // Can be null if in progress
}

// BackupCreateOptions defines the optional request body for creating a new backup.
type BackupCreateOptions struct {
	// A name for the backup. If nil, a name will be generated based on the timestamp.
	Name *string `json:"name,omitempty"`

	// A string containing a list of files and folders to ignore, with each entry
	// on a new line. If nil, the server's .pteroignore file will be used.
	Ignored *string `json:"ignored,omitempty"`
}

// BackupDownload contains the signed URL for downloading a backup.
type BackupDownload struct {
	URL string `json:"url"`
}

// backupResponse is a helper for unmarshaling single backup responses.
type BackupResponse struct {
	Object     string  `json:"object"`
	Attributes *Backup `json:"attributes"`
}

// backupDownloadResponse is a helper for unmarshaling the download link response.
type BackupDownloadResponse struct {
	Object     string          `json:"object"`
	Attributes *BackupDownload `json:"attributes"`
}
