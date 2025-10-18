package application_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/application"
)

func TestNodes_ListNodes(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nodes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "node",
				"attributes": {
					"id": 1,
					"name": "Node 1",
					"fqdn": "node1.example.com",
					"memory": 8192,
					"disk": 102400
				}
			}],
			"meta": {"pagination": {"total": 1, "total_pages": 1}}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	nodes, _, err := appClient.ListNodes(context.Background(), pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListNodes() error = %v", err)
	}

	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}
}

func TestNodes_GetNodeConfiguration(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nodes/1/configuration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"debug": false,
			"uuid": "test-uuid",
			"token_id": "test-token-id",
			"token": "test-token",
			"api": {
				"host": "0.0.0.0",
				"port": 8080,
				"ssl": {
					"enabled": true,
					"cert": "/etc/letsencrypt/live/node.example.com/fullchain.pem",
					"key": "/etc/letsencrypt/live/node.example.com/privkey.pem"
				},
				"upload_limit": 100
			},
			"system": {
				"root_directory": "/var/lib/pterodactyl/volumes",
				"log_directory": "/var/log/pterodactyl",
				"data": "/var/lib/pterodactyl",
				"sftp": {
					"bind_port": 2022
				},
				"crash_detection": {
					"enabled": true,
					"timeout": 60
				},
				"backups": {
					"write_limit": 0
				},
				"transfers": {
					"download_limit": 0
				}
			},
			"allowed_mounts": [],
			"remote": "https://panel.example.com"
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	config, err := appClient.GetNodeConfiguration(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNodeConfiguration() error = %v", err)
	}

	if config.UUID != "test-uuid" {
		t.Errorf("UUID = %s, want test-uuid", config.UUID)
	}
	if config.API.Port != 8080 {
		t.Errorf("API.Port = %d, want 8080", config.API.Port)
	}
}

func TestNodes_GetDeployableNodes(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nodes/deployable", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		memory := r.URL.Query().Get("memory")
		disk := r.URL.Query().Get("disk")

		if memory != "1024" || disk != "5120" {
			t.Errorf("Query params: memory=%s, disk=%s; want memory=1024, disk=5120", memory, disk)
		}

		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "node",
				"attributes": {
					"id": 1,
					"name": "Available Node"
				}
			}]
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	nodes, err := appClient.GetDeployableNodes(context.Background(), 1024, 5120)
	if err != nil {
		t.Fatalf("GetDeployableNodes() error = %v", err)
	}

	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}
}
