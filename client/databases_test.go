package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestDatabases_CreateDatabase(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/databases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Decode() failed: %v", err)
		}
		if req["database"] != "testdb" || req["remote"] != "%" {
			t.Errorf("unexpected request body: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"server_database","attributes":{"id":"s1_1","name":"s1_testdb"}}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	db, err := clientAPI.CreateDatabase(context.Background(), "d3aac109", "testdb", "%")
	if err != nil {
		t.Fatalf("CreateDatabase() returned error: %v", err)
	}
	if db.ID != "s1_1" {
		t.Errorf("expected database ID s1_1, got %s", db.ID)
	}
}
