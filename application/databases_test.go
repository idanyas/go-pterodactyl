package application_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/application"
)

func TestDatabases_ListServerDatabases(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/servers/1/databases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "server_database",
				"attributes": {
					"id": 1,
					"server": 1,
					"host": 1,
					"database": "s1_test",
					"username": "u1_test",
					"remote": "%"
				}
			}],
			"meta": {"pagination": {"total": 1, "total_pages": 1}}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	databases, _, err := appClient.ListServerDatabases(context.Background(), 1, pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListServerDatabases() error = %v", err)
	}

	if len(databases) != 1 {
		t.Errorf("expected 1 database, got %d", len(databases))
	}
}

func TestDatabases_CreateServerDatabase(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	req := application.CreateServerDatabaseRequest{
		Database: "testdb",
		Remote:   "%",
		Host:     1,
	}

	mux.HandleFunc("/api/application/servers/1/databases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var received application.CreateServerDatabaseRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if received.Database != req.Database {
			t.Errorf("Database = %s, want %s", received.Database, req.Database)
		}

		fmt.Fprint(w, `{
			"object": "server_database",
			"attributes": {
				"id": 2,
				"database": "s1_testdb"
			}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	db, err := appClient.CreateServerDatabase(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("CreateServerDatabase() error = %v", err)
	}

	if db.ID != 2 {
		t.Errorf("ID = %d, want 2", db.ID)
	}
}

func TestDatabases_ResetServerDatabasePassword(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/servers/1/databases/2/reset-password", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	err := appClient.ResetServerDatabasePassword(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("ResetServerDatabasePassword() error = %v", err)
	}
}

func TestDatabases_DeleteServerDatabase(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/servers/1/databases/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	err := appClient.DeleteServerDatabase(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("DeleteServerDatabase() error = %v", err)
	}
}
