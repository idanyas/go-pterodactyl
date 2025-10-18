// Package pterodactyl provides a comprehensive Go client for the Pterodactyl Panel API.
package pterodactyl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/idanyas/go-pterodactyl/application"
	"github.com/idanyas/go-pterodactyl/client"
	"github.com/idanyas/go-pterodactyl/pagination"
	"github.com/idanyas/go-pterodactyl/transport"
)

const (
	// APIVersion is the version of the Pterodactyl API this client targets.
	APIVersion = "v1"
	// defaultUserAgent is the default User-Agent header sent with requests.
	defaultUserAgent = "go-pterodactyl/v1.0"
)

// ListOptions specifies optional parameters to list methods.
// It is an alias for pagination.ListOptions.
type ListOptions = pagination.ListOptions

// Client is the primary client for interacting with the Pterodactyl API.
// It provides access to the Application and Client APIs.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client

	// API Clients
	app    application.ApplicationClient
	client client.ClientClient
}

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithAPIKey sets the API key to be used for authentication.
// The key should be prefixed with `ptla_` for the Application API
// or `ptlc_` for the Client API.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithHTTPClient sets a custom http.Client for the Pterodactyl client.
// This is useful for configuring custom transports, timeouts, or other settings.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// New creates a new Pterodactyl API client.
//
// panelURL is the base URL of the Pterodactyl panel (e.g., "https://panel.example.com").
// opts are functional options to configure the client, such as setting the API key.
func New(panelURL string, opts ...Option) (*Client, error) {
	if panelURL == "" {
		return nil, fmt.Errorf("panelURL cannot be empty")
	}

	baseURL, err := url.Parse(strings.TrimRight(panelURL, "/") + "/api/")
	if err != nil {
		return nil, fmt.Errorf("failed to parse panelURL: %w", err)
	}

	c := &Client{
		baseURL: baseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	// Decorate the transport with authentication, rate limiting, and retries.
	baseTransport := c.httpClient.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	c.httpClient.Transport = transport.New(
		baseTransport,
		c.apiKey,
		APIVersion,
		defaultUserAgent,
	)

	c.app = application.New(c)
	c.client = client.New(c)

	return c, nil
}

// Application returns a client for interacting with the Application API.
func (c *Client) Application() application.ApplicationClient {
	return c.app
}

// Client returns a client for interacting with the Client API.
func (c *Client) Client() client.ClientClient {
	return c.client
}

// newRequest creates an API request. A relative URL path can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	var r io.ReadWriter
	if body != nil {
		r = new(bytes.Buffer)
		if err := json.NewEncoder(r).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), r)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response bodies
		}
	}

	return resp, err
}

// Do performs a request. It is the underlying method for all API calls.
func (c *Client) Do(ctx context.Context, method, path string, body, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

// DoRequest performs a raw request, allowing for more control.
func (c *Client) DoRequest(req *http.Request, v interface{}) (*http.Response, error) {
	return c.do(req, v)
}
