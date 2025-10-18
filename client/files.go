package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// ChmodFileRequest represents a file permission change request.
type ChmodFileRequest struct {
	File string `json:"file"`
	Mode string `json:"mode"`
}

// ListFiles retrieves the contents of a server directory.
func (c *client) ListFiles(ctx context.Context, serverID, directory string) ([]*models.FileObject, error) {
	path := fmt.Sprintf("client/servers/%s/files/list?directory=%s", serverID, directory)
	var response struct {
		Data []struct {
			Attributes models.FileObject `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	files := make([]*models.FileObject, len(response.Data))
	for i, item := range response.Data {
		files[i] = &item.Attributes
	}
	return files, nil
}

// GetFileContents retrieves the contents of a specific file.
func (c *client) GetFileContents(ctx context.Context, serverID, filePath string) (string, error) {
	path := fmt.Sprintf("client/servers/%s/files/contents?file=%s", serverID, filePath)
	var content string
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &content)
	return content, err
}

// WriteFile creates or updates a file with new content.
func (c *client) WriteFile(ctx context.Context, serverID, filePath, content string) error {
	path := fmt.Sprintf("client/servers/%s/files/write?file=%s", serverID, filePath)
	_, err := c.client.Do(ctx, http.MethodPost, path, content, nil)
	return err
}

// CreateDirectory creates a new directory on the server.
func (c *client) CreateDirectory(ctx context.Context, serverID, root, name string) error {
	path := fmt.Sprintf("client/servers/%s/files/create-folder", serverID)
	req := map[string]string{"root": root, "name": name}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// DeleteFiles permanently deletes files or directories.
func (c *client) DeleteFiles(ctx context.Context, serverID, root string, files []string) error {
	path := fmt.Sprintf("client/servers/%s/files/delete", serverID)
	req := map[string]interface{}{"root": root, "files": files}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// RenameFile renames or moves a file or directory.
func (c *client) RenameFile(ctx context.Context, serverID, root, from, to string) error {
	path := fmt.Sprintf("client/servers/%s/files/rename", serverID)
	req := map[string]interface{}{
		"root": root,
		"files": []map[string]string{
			{"from": from, "to": to},
		},
	}
	_, err := c.client.Do(ctx, http.MethodPut, path, req, nil)
	return err
}

// CopyFile creates a copy of a file or directory.
func (c *client) CopyFile(ctx context.Context, serverID, location string) error {
	path := fmt.Sprintf("client/servers/%s/files/copy", serverID)
	req := map[string]string{"location": location}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// GetDownloadURL retrieves a pre-signed URL for downloading a file.
func (c *client) GetDownloadURL(ctx context.Context, serverID, filePath string) (*models.SignedURL, error) {
	path := fmt.Sprintf("client/servers/%s/files/download?file=%s", serverID, filePath)
	var response struct {
		Attributes models.SignedURL `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// GetUploadURL retrieves a pre-signed URL for uploading files.
func (c *client) GetUploadURL(ctx context.Context, serverID, directory string) (*models.SignedURL, error) {
	path := fmt.Sprintf("client/servers/%s/files/upload?directory=%s", serverID, directory)
	var response struct {
		Attributes models.SignedURL `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// CompressFiles creates an archive from specified files and directories.
func (c *client) CompressFiles(ctx context.Context, serverID, root string, files []string) (*models.FileObject, error) {
	path := fmt.Sprintf("client/servers/%s/files/compress", serverID)
	req := map[string]interface{}{"root": root, "files": files}
	var response struct {
		Attributes models.FileObject `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DecompressFile extracts files from an archive.
func (c *client) DecompressFile(ctx context.Context, serverID, root, file string) error {
	path := fmt.Sprintf("client/servers/%s/files/decompress", serverID)
	req := map[string]string{"root": root, "file": file}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// ChmodFiles modifies file or directory permissions.
func (c *client) ChmodFiles(ctx context.Context, serverID, root string, files []ChmodFileRequest) error {
	path := fmt.Sprintf("client/servers/%s/files/chmod", serverID)
	req := map[string]interface{}{"root": root, "files": files}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}

// PullFile downloads a file from a URL directly to the server.
func (c *client) PullFile(ctx context.Context, serverID, url, directory, filename string) error {
	path := fmt.Sprintf("client/servers/%s/files/pull", serverID)
	req := map[string]string{
		"url":       url,
		"directory": directory,
		"filename":  filename,
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, nil)
	return err
}
