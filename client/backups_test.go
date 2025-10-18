package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestBackups_ListBackups(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/backups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "backup",
				"attributes": {
					"uuid": "backup-uuid",
					"name": "Test Backup",
					"bytes": 1024,
					"is_successful": true,
					"is_locked": false
				}
			}]
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	backups, err := clientAPI.ListBackups(context.Background(), "d3aac109")
	if err != nil {
		t.Fatalf("ListBackups() error = %v", err)
	}

	if len(backups) != 1 {
		t.Errorf("expected 1 backup, got %d", len(backups))
	}
}

func TestBackups_CreateBackup(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	req := client.CreateBackupRequest{
		Name:     "New Backup",
		Ignored:  "*.log",
		IsLocked: true,
	}

	mux.HandleFunc("/api/client/servers/d3aac109/backups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var received client.CreateBackupRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if received.Name != req.Name {
			t.Errorf("Name = %s, want %s", received.Name, req.Name)
		}

		fmt.Fprint(w, `{
			"object": "backup",
			"attributes": {
				"uuid": "new-backup-uuid",
				"name": "New Backup",
				"is_locked": true
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	backup, err := clientAPI.CreateBackup(context.Background(), "d3aac109", req)
	if err != nil {
		t.Fatalf("CreateBackup() error = %v", err)
	}

	if backup.Name != req.Name {
		t.Errorf("Name = %s, want %s", backup.Name, req.Name)
	}
}

func TestBackups_GetBackupDownloadURL(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/backups/backup-uuid/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "signed_url",
			"attributes": {
				"url": "https://example.com/download/backup.tar.gz"
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	signedURL, err := clientAPI.GetBackupDownloadURL(context.Background(), "d3aac109", "backup-uuid")
	if err != nil {
		t.Fatalf("GetBackupDownloadURL() error = %v", err)
	}

	if signedURL.URL == "" {
		t.Error("expected non-empty URL")
	}
}

func TestBackups_DeleteBackup(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/backups/backup-uuid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.DeleteBackup(context.Background(), "d3aac109", "backup-uuid")
	if err != nil {
		t.Fatalf("DeleteBackup() error = %v", err)
	}
}
