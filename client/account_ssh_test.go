package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestSSH_ListSSHKeys(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/account/ssh-keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "ssh_key",
				"attributes": {
					"name": "My Key",
					"fingerprint": "SHA256:abcd1234",
					"public_key": "ssh-rsa AAAAB3..."
				}
			}]
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	keys, err := clientAPI.ListSSHKeys(context.Background())
	if err != nil {
		t.Fatalf("ListSSHKeys() error = %v", err)
	}

	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
}

func TestSSH_AddSSHKey(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/account/ssh-keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if req["name"] != "New Key" {
			t.Errorf("name = %s, want New Key", req["name"])
		}

		fmt.Fprint(w, `{
			"object": "ssh_key",
			"attributes": {
				"name": "New Key",
				"fingerprint": "SHA256:newkey1234"
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	key, err := clientAPI.AddSSHKey(context.Background(), "New Key", "ssh-rsa AAAAB3...")
	if err != nil {
		t.Fatalf("AddSSHKey() error = %v", err)
	}

	if key.Name != "New Key" {
		t.Errorf("Name = %s, want New Key", key.Name)
	}
}

func TestSSH_RemoveSSHKey(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/account/ssh-keys/remove", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.RemoveSSHKey(context.Background(), "SHA256:abcd1234")
	if err != nil {
		t.Fatalf("RemoveSSHKey() error = %v", err)
	}
}
