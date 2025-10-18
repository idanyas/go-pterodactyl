package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestNetwork_ListAllocations(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/network/allocations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "allocation",
				"attributes": {
					"id": 1,
					"ip": "192.168.1.100",
					"port": 25565,
					"is_default": true
				}
			}]
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	allocations, err := clientAPI.ListAllocations(context.Background(), "d3aac109")
	if err != nil {
		t.Fatalf("ListAllocations() error = %v", err)
	}

	if len(allocations) != 1 {
		t.Errorf("expected 1 allocation, got %d", len(allocations))
	}
}

func TestNetwork_SetPrimaryAllocation(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/network/allocations/2/primary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.SetPrimaryAllocation(context.Background(), "d3aac109", 2)
	if err != nil {
		t.Fatalf("SetPrimaryAllocation() error = %v", err)
	}
}

func TestNetwork_DeleteAllocation(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/network/allocations/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.DeleteAllocation(context.Background(), "d3aac109", 2)
	if err != nil {
		t.Fatalf("DeleteAllocation() error = %v", err)
	}
}
