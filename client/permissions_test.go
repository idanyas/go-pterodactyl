package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestPermissions_GetSystemPermissions(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/permissions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "system_permissions",
			"attributes": {
				"permissions": {
					"control": {
						"description": "Power control permissions",
						"keys": {
							"console": "Send console commands",
							"start": "Start server",
							"stop": "Stop server"
						}
					}
				}
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	perms, err := clientAPI.GetSystemPermissions(context.Background())
	if err != nil {
		t.Fatalf("GetSystemPermissions() error = %v", err)
	}

	if len(perms.Permissions) == 0 {
		t.Error("expected permissions, got none")
	}

	if control, ok := perms.Permissions["control"]; ok {
		if control.Description == "" {
			t.Error("expected non-empty description")
		}
		if len(control.Keys) == 0 {
			t.Error("expected permission keys")
		}
	} else {
		t.Error("expected 'control' permission group")
	}
}
