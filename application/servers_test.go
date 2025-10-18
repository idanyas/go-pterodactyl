package application_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/application"
)

func TestServers_SuspendServer(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/servers/1/suspend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	err := appClient.SuspendServer(context.Background(), 1)
	if err != nil {
		t.Fatalf("SuspendServer returned error: %v", err)
	}
}

func TestServers_DeleteServer(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/servers/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	err := appClient.DeleteServer(context.Background(), 1, false)
	if err != nil {
		t.Fatalf("DeleteServer returned error: %v", err)
	}
}

func TestServers_DeleteServer_Force(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/servers/1/force", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	err := appClient.DeleteServer(context.Background(), 1, true)
	if err != nil {
		t.Fatalf("DeleteServer(force) returned error: %v", err)
	}
}
