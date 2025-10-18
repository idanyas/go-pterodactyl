package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestSettings_RenameServer(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/settings/rename", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if req["name"] != "New Name" {
			t.Errorf("name = %s, want New Name", req["name"])
		}

		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.RenameServer(context.Background(), "d3aac109", "New Name", "New Description")
	if err != nil {
		t.Fatalf("RenameServer() error = %v", err)
	}
}

func TestSettings_ReinstallServer(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/settings/reinstall", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusAccepted)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.ReinstallServer(context.Background(), "d3aac109")
	if err != nil {
		t.Fatalf("ReinstallServer() error = %v", err)
	}
}

func TestSettings_UpdateDockerImage(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/settings/docker-image", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if req["docker_image"] != "ghcr.io/pterodactyl/yolks:java_17" {
			t.Errorf("docker_image = %s, want ghcr.io/pterodactyl/yolks:java_17", req["docker_image"])
		}

		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.UpdateDockerImage(context.Background(), "d3aac109", "ghcr.io/pterodactyl/yolks:java_17")
	if err != nil {
		t.Fatalf("UpdateDockerImage() error = %v", err)
	}
}
