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

func TestLocations_ListLocations(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/locations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"object": "list",
			"data": [
				{
					"object": "location",
					"attributes": {
						"id": 1,
						"short": "us-west",
						"long": "United States - West"
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

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	locations, _, err := appClient.ListLocations(context.Background(), pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListLocations() error = %v", err)
	}

	if len(locations) != 1 {
		t.Errorf("expected 1 location, got %d", len(locations))
	}
	if locations[0].Short != "us-west" {
		t.Errorf("expected short 'us-west', got %s", locations[0].Short)
	}
}

func TestLocations_CreateLocation(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	req := application.CreateLocationRequest{
		Short: "eu-central",
		Long:  "Europe - Central",
	}

	mux.HandleFunc("/api/application/locations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var received application.CreateLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if received.Short != req.Short {
			t.Errorf("Short = %s, want %s", received.Short, req.Short)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"object": "location",
			"attributes": {
				"id": 2,
				"short": "eu-central",
				"long": "Europe - Central"
			}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	location, err := appClient.CreateLocation(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateLocation() error = %v", err)
	}

	if location.Short != req.Short {
		t.Errorf("Short = %s, want %s", location.Short, req.Short)
	}
}

func TestLocations_GetLocation(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/locations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "location",
			"attributes": {
				"id": 1,
				"short": "us-west",
				"long": "United States - West"
			}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	location, err := appClient.GetLocation(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLocation() error = %v", err)
	}

	if location.ID != 1 {
		t.Errorf("ID = %d, want 1", location.ID)
	}
}

func TestLocations_UpdateLocation(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	req := application.UpdateLocationRequest{
		Long: "United States - West Coast",
	}

	mux.HandleFunc("/api/application/locations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		fmt.Fprint(w, `{
			"object": "location",
			"attributes": {
				"id": 1,
				"short": "us-west",
				"long": "United States - West Coast"
			}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	location, err := appClient.UpdateLocation(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("UpdateLocation() error = %v", err)
	}

	if location.Long != req.Long {
		t.Errorf("Long = %s, want %s", location.Long, req.Long)
	}
}

func TestLocations_DeleteLocation(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/locations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	err := appClient.DeleteLocation(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteLocation() error = %v", err)
	}
}

func TestLocations_InvalidID(t *testing.T) {
	_, serverURL, teardown := setup()
	defer teardown()

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "GetLocation with 0",
			fn: func() error {
				_, err := appClient.GetLocation(context.Background(), 0)
				return err
			},
		},
		{
			name: "GetLocation with negative",
			fn: func() error {
				_, err := appClient.GetLocation(context.Background(), -1)
				return err
			},
		},
		{
			name: "DeleteLocation with 0",
			fn: func() error {
				return appClient.DeleteLocation(context.Background(), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Error("expected error for invalid ID, got nil")
			}
		})
	}
}
