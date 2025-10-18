package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/client"
)

func TestActivity_ListAccountActivity(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/account/activity", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "activity_log",
				"attributes": {
					"id": "log-uuid",
					"event": "user:account.email-changed",
					"is_api": false,
					"ip": "192.168.1.1",
					"description": "Email changed"
				}
			}],
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

	logs, _, err := clientAPI.ListAccountActivity(context.Background(), pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListAccountActivity() error = %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
}

func TestActivity_ListServerActivity(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/activity", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "activity_log",
				"attributes": {
					"id": "log-uuid",
					"event": "server:power.start",
					"is_api": false,
					"ip": "192.168.1.1",
					"description": "Server started"
				}
			}],
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

	logs, _, err := clientAPI.ListServerActivity(context.Background(), "d3aac109", pterodactyl.ListOptions{})
	if err != nil {
		t.Fatalf("ListServerActivity() error = %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
}
