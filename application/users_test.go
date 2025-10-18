package application_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/application"
	"github.com/idanyas/go-pterodactyl/pagination"
)

func TestUsers_ListUsers(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/application/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"object": "list",
			"data": [
				{
					"object": "user",
					"attributes": { "id": 1, "username": "testuser" }
				}
			],
			"meta": { "pagination": { "total": 1, "total_pages": 1 } }
		}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	users, _, err := appClient.ListUsers(context.Background(), pagination.ListOptions{})
	if err != nil {
		t.Fatalf("ListUsers returned error: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
	if users[0].ID != 1 || users[0].Username != "testuser" {
		t.Errorf("unexpected user data: %+v", users[0])
	}
}

func TestUsers_CreateUser(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	req := application.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	mux.HandleFunc("/api/application/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var v application.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("Decode() failed: %v", err)
		}
		if !reflect.DeepEqual(v, req) {
			t.Errorf("Request body = %+v, want %+v", v, req)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"user","attributes":{"id":1, "username":"testuser"}}`)
	})

	client, _ := pterodactyl.New(serverURL, pterodactyl.WithAPIKey("test-key"))
	appClient := application.New(client)

	user, err := appClient.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}
