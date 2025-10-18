package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestAccount_GetAccount(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/account", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", "Application/vnd.pterodactyl.v1+json")
		testHeader(t, r, "Authorization", "Bearer test-key")
		fmt.Fprint(w, `{"object":"user","attributes":{"id":1,"username":"testuser","email":"test@example.com"}}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	user, err := clientAPI.GetAccount(context.Background())
	if err != nil {
		t.Fatalf("GetAccount() returned error: %v", err)
	}

	if user.ID != 1 || user.Username != "testuser" {
		t.Errorf("unexpected user data: %+v", user)
	}
}
