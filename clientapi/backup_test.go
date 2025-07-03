package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/errors"
	"github.com/davidarkless/go-pterodactyl/internal/testutil"
)

const testServerIdentifier = "test-server"

func TestBackupsService_List(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	checksum := "test-checksum"
	expectedBackups := []*api.Backup{
		{
			UUID:         "uuid1",
			IsSuccessful: true,
			IsLocked:     false,
			Name:         "backup1",
			Bytes:        1024,
			CreatedAt:    now,
			CompletedAt:  &now,
		},
		{
			UUID:         "uuid2",
			IsSuccessful: true,
			IsLocked:     true,
			Name:         "backup2",
			IgnoredFiles: []string{"/ignore.txt"},
			Checksum:     &checksum,
			Bytes:        2048,
			CreatedAt:    now,
			CompletedAt:  &now,
		},
	}

	data := make([]*api.ListItem[api.Backup], len(expectedBackups))
	for i, backup := range expectedBackups {
		data[i] = &api.ListItem[api.Backup]{Object: "backup", Attributes: backup}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.Backup]{
		Object: "list",
		Data:   data,
		Meta:   meta,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		backups, m, err := s.List(context.Background(), api.PaginationOptions{Page: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(backups, expectedBackups) {
			t.Errorf("expected backups %+v, got %+v", expectedBackups, backups)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, req.Method)
		}
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusInternalServerError,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusInternalServerError},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		_, _, err := s.List(context.Background(), api.PaginationOptions{})
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Create(t *testing.T) {
	name := "new-backup"
	options := api.BackupCreateOptions{Name: &name}
	jsonOptions, _ := json.Marshal(options)

	expectedBackup := &api.Backup{
		UUID:         "new-uuid",
		IsSuccessful: false, // In progress
		Name:         name,
		CreatedAt:    time.Now().Truncate(time.Second),
	}
	res := api.BackupResponse{
		Object:     "backup",
		Attributes: expectedBackup,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK, // The API returns 200 on create
				Body:       jsonBody,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		backup, err := s.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(backup, expectedBackup) {
			t.Errorf("expected backup %+v, got %+v", expectedBackup, backup)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, req.Method)
		}
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusLocked,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusLocked},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		_, err := s.Create(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Details(t *testing.T) {
	uuid := "test-uuid"
	expectedBackup := &api.Backup{
		UUID:         uuid,
		IsSuccessful: true,
		Name:         "details-backup",
		Bytes:        4096,
		CreatedAt:    time.Now().Truncate(time.Second),
	}
	res := api.BackupResponse{
		Object:     "backup",
		Attributes: expectedBackup,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		backup, err := s.Details(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(backup, expectedBackup) {
			t.Errorf("expected backup %+v, got %+v", expectedBackup, backup)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		_, err := s.Details(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Download(t *testing.T) {
	uuid := "test-uuid"
	expectedDownload := &api.BackupDownload{URL: "https://example.com/download/backup.zip"}
	res := api.BackupDownloadResponse{
		Object:     "backup_download",
		Attributes: expectedDownload,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		download, err := s.Download(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(download, expectedDownload) {
			t.Errorf("expected download %+v, got %+v", expectedDownload, download)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s/download", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		_, err := s.Download(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Delete(t *testing.T) {
	uuid := "test-uuid-delete"

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
