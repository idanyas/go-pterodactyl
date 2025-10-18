package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/client"
)

func TestServers_ListServers(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"object": "list",
			"data": [
				{
					"object": "server",
					"attributes": {
						"identifier": "d3aac109",
						"uuid": "d3aac109-e5e0-4331-b03e-3454f7e136dc",
						"name": "Test Server"
					}
				}
			],
			"meta": {
				"pagination": {
					"total": 1,
					"total_pages": 1
				}
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	servers, _, err := clientAPI.ListServers(context.Background(), pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListServers() returned error: %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(servers))
	}
	if servers[0].Identifier != "d3aac109" {
		t.Errorf("unexpected server identifier: %s", servers[0].Identifier)
	}
}

func TestServers_SendPowerAction(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/power", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Decode() failed: %v", err)
		}
		if req["signal"] != "restart" {
			t.Errorf("signal = %s, want restart", req["signal"])
		}

		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.SendPowerAction(context.Background(), "d3aac109", "restart")
	if err != nil {
		t.Fatalf("SendPowerAction() returned error: %v", err)
	}
}
