package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestSubusers_ListSubusers(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "subuser",
				"attributes": {
					"uuid": "user-uuid",
					"username": "john",
					"email": "john@example.com",
					"permissions": ["control.start", "control.stop"]
				}
			}]
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	subusers, err := clientAPI.ListSubusers(context.Background(), "d3aac109")
	if err != nil {
		t.Fatalf("ListSubusers() error = %v", err)
	}

	if len(subusers) != 1 {
		t.Errorf("expected 1 subuser, got %d", len(subusers))
	}
	if subusers[0].Username != "john" {
		t.Errorf("Username = %s, want john", subusers[0].Username)
	}
}

func TestSubusers_CreateSubuser(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if req["email"] != "newuser@example.com" {
			t.Errorf("email = %v, want newuser@example.com", req["email"])
		}

		fmt.Fprint(w, `{
			"object": "subuser",
			"attributes": {
				"uuid": "new-user-uuid",
				"email": "newuser@example.com"
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	subuser, err := clientAPI.CreateSubuser(context.Background(), "d3aac109", "newuser@example.com", []string{"control.console"})
	if err != nil {
		t.Fatalf("CreateSubuser() error = %v", err)
	}

	if subuser.Email != "newuser@example.com" {
		t.Errorf("Email = %s, want newuser@example.com", subuser.Email)
	}
}

func TestSubusers_DeleteSubuser(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/users/user-uuid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.DeleteSubuser(context.Background(), "d3aac109", "user-uuid")
	if err != nil {
		t.Fatalf("DeleteSubuser() error = %v", err)
	}
}
