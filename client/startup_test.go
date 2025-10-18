package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestStartup_GetStartupConfig(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/startup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "egg_variable",
				"attributes": {
					"name": "Server Jar",
					"env_variable": "SERVER_JARFILE",
					"default_value": "server.jar",
					"server_value": "server.jar",
					"is_editable": true
				}
			}],
			"meta": {
				"startup_command": "java -jar {{SERVER_JARFILE}}",
				"raw_startup_command": "java -jar server.jar",
				"docker_images": {
					"ghcr.io/pterodactyl/yolks:java_17": "Java 17"
				}
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	config, err := clientAPI.GetStartupConfig(context.Background(), "d3aac109")
	if err != nil {
		t.Fatalf("GetStartupConfig() error = %v", err)
	}

	if len(config.Variables) != 1 {
		t.Errorf("expected 1 variable, got %d", len(config.Variables))
	}
	if config.StartupCommand == "" {
		t.Error("expected non-empty startup command")
	}
}

func TestStartup_UpdateStartupVariable(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/startup/variable", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if req["key"] != "SERVER_JARFILE" {
			t.Errorf("key = %s, want SERVER_JARFILE", req["key"])
		}
		if req["value"] != "custom.jar" {
			t.Errorf("value = %s, want custom.jar", req["value"])
		}

		fmt.Fprint(w, `{
			"object": "egg_variable",
			"attributes": {
				"server_value": "custom.jar"
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	variable, err := clientAPI.UpdateStartupVariable(context.Background(), "d3aac109", "SERVER_JARFILE", "custom.jar")
	if err != nil {
		t.Fatalf("UpdateStartupVariable() error = %v", err)
	}

	if variable.ServerValue != "custom.jar" {
		t.Errorf("ServerValue = %s, want custom.jar", variable.ServerValue)
	}
}
