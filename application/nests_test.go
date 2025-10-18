package application_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/application"
)

func TestNests_ListNests(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "nest",
				"attributes": {
					"id": 1,
					"name": "Minecraft",
					"description": "Minecraft servers"
				}
			}],
			"meta": {"pagination": {"total": 1, "total_pages": 1}}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	nests, _, err := appClient.ListNests(context.Background(), pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListNests() error = %v", err)
	}

	if len(nests) != 1 {
		t.Errorf("expected 1 nest, got %d", len(nests))
	}
	if nests[0].Name != "Minecraft" {
		t.Errorf("Name = %s, want Minecraft", nests[0].Name)
	}
}

func TestNests_GetNest(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nests/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "nest",
			"attributes": {
				"id": 1,
				"name": "Minecraft"
			}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	nest, err := appClient.GetNest(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNest() error = %v", err)
	}

	if nest.ID != 1 {
		t.Errorf("ID = %d, want 1", nest.ID)
	}
}

func TestNests_ListNestEggs(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nests/1/eggs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "egg",
				"attributes": {
					"id": 1,
					"name": "Vanilla",
					"nest": 1
				}
			}],
			"meta": {"pagination": {"total": 1, "total_pages": 1}}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	eggs, _, err := appClient.ListNestEggs(context.Background(), 1, pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListNestEggs() error = %v", err)
	}

	if len(eggs) != 1 {
		t.Errorf("expected 1 egg, got %d", len(eggs))
	}
}

func TestNests_GetEgg(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/nests/1/eggs/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "egg",
			"attributes": {
				"id": 1,
				"name": "Vanilla",
				"nest": 1
			}
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	egg, err := appClient.GetEgg(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("GetEgg() error = %v", err)
	}

	if egg.Name != "Vanilla" {
		t.Errorf("Name = %s, want Vanilla", egg.Name)
	}
}
