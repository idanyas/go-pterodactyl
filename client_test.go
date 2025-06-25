package pterodactyl_test

import (
	"github.com/davidarkless/go-pterodactyl"
	"testing"
)

func TestNewClientSuccess(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name          string
		apiKey        string
		keyType       pterodactyl.KeyType
		expectedError bool
	}{
		{
			name:          "Valid Application Key",
			apiKey:        "ptla_abc123", // A dummy key with the correct prefix
			keyType:       pterodactyl.ApplicationKey,
			expectedError: false,
		},
		{
			name:          "Valid Client Key",
			apiKey:        "ptlc_def456", // A dummy key with the correct prefix
			keyType:       pterodactyl.ClientKey,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := pterodactyl.NewClient("https://fake-panel.com", tc.apiKey, tc.keyType)

			if tc.expectedError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}

				// On success, we should also check that the client is not nil and its services are initialized.
				if client == nil {
					t.Fatal("expected client to be non-nil on success")
				}
				if client.Application == nil {
					t.Error("expected Application service to be initialized")
				}
				if client.Client == nil {
					t.Error("expected Client service to be initialized")
				}
			}
		})
	}
}

// TestNewClient_InvalidKeyFormat checks for errors when API keys have incorrect prefixes.
func TestNewClient_InvalidKeyFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		apiKey        string
		keyType       pterodactyl.KeyType
		expectedError bool
	}{
		{
			name:          "Invalid Application Key",
			apiKey:        "ptlc_wrongprefix", // A client key prefix used for an application client
			keyType:       pterodactyl.ApplicationKey,
			expectedError: true,
		},
		{
			name:          "Invalid Client Key",
			apiKey:        "ptla_wrongprefix", // An application key prefix used for a client client
			keyType:       pterodactyl.ClientKey,
			expectedError: true,
		},
		{
			name:          "Application Key without prefix",
			apiKey:        "noprefix",
			keyType:       pterodactyl.ApplicationKey,
			expectedError: true,
		},
		{
			name:          "Client Key without prefix",
			apiKey:        "noprefix",
			keyType:       pterodactyl.ClientKey,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := pterodactyl.NewClient("https://fake-panel.com", tc.apiKey, tc.keyType)

			if !tc.expectedError {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected an error for invalid key format but got none")
				}
			}
		})
	}
}

// TestNewClient_InvalidURL demonstrates testing for a malformed base URL.
// While this case might be less common, it's good practice to ensure robustness.
func TestNewClient_InvalidURL(t *testing.T) {
	t.Parallel()

	// This is a special character that will cause url.Parse to fail in NewRequest
	malformedURL := "::not a valid url"

	client, err := pterodactyl.NewClient(malformedURL, "ptlc_dummykey", pterodactyl.ClientKey)
	if err != nil {
		// The error doesn't happen in NewClient itself, but on the first request.
		// Let's test that by trying to make a request.
		// Since we can't access client.NewRequest directly (it's unexported in our test package),
		// we test a public method that uses it.
		// We'll need to update the client.go NewClient to check the URL at creation time.

		// Let's go back to client.go and improve it first.
		t.Skip("Skipping test: NewClient should validate the baseURL upon creation.")
	}

	// This test reveals a small design flaw: the baseURL isn't validated until a request is made.
	// Let's fix that in pterodactyl.go first, then complete this test.
	// See the "Refinement" section below.

	// For now, this is how you'd test the *current* behavior:
	// Assuming `client.Client.ListPermissions()` is a simple method with no body.
	_, listErr := client.Client.ListPermissions()
	if listErr == nil {
		t.Errorf("expected an error from a request with a malformed baseURL, but got none")
	}
}
