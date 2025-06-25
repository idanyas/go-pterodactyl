package pterodactyl

import (
	"encoding/json"
	"fmt"
	"github.com/davidarkless/go-pterodactyl/api"
	"github.com/davidarkless/go-pterodactyl/application"
	"github.com/davidarkless/go-pterodactyl/clientapi"
	"github.com/davidarkless/go-pterodactyl/errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type KeyType int

const (
	ApplicationKey KeyType = iota
	ClientKey
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client

	Application *application.Application
	Client      *clientapi.ClientAPIService
}

func NewClient(baseURL, apiKey string, keyType KeyType) (*Client, error) {

	if keyType == ApplicationKey && !strings.HasPrefix(apiKey, "ptla_") {
		return nil, fmt.Errorf("invalid application key: key must start with 'ptla_'")
	}
	if keyType == ClientKey && !strings.HasPrefix(apiKey, "ptlc_") {
		return nil, fmt.Errorf("invalid client key: key must start with 'ptlc_'")
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}
	client := &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}

	client.Application = &application.Application{}
	client.Application.Users = application.NewUsersService(client)
	client.Application.Nodes = application.NewNodesService(client)
	client.Application.Locations = application.NewLocationService(client)
	client.Application.Servers = application.NewServersService(client)
	client.Application.Nests = application.NewNestsService(client)

	client.Client = clientapi.NewClientAPI(client)

	return client, nil
}

func (c *Client) NewRequest(method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error) {

	rel, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if options != nil {
		q := rel.Query()
		if options.Page > 0 {
			q.Set("page", strconv.Itoa(options.Page))
		}
		if options.PerPage > 0 {
			q.Set("per_page", strconv.Itoa(options.PerPage))
		}
		if len(options.Include) > 0 {
			q.Set("include", strings.Join(options.Include, ","))
		}
		rel.RawQuery = q.Encode()
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	fullURL := u.ResolveReference(rel)

	req, err := http.NewRequest(method, fullURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close() // ignore error

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// Error handling logic
		apiErr := &errors.APIError{HTTPStatusCode: res.StatusCode}
		if err = json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("pterodactyl: API error (status %d), failed to parse error response: %w", res.StatusCode, err)
		}
		return nil, apiErr
	}

	// If v is not nil, decode the successful response body into it.
	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return nil, fmt.Errorf("failed to decode successful response: %w", err)
		}
	}

	return res, nil
}

// unmarshalList is an internal helper that decodes a paginated list response
// from the Pterodactyl API and flattens it into a simple slice of models.
// It uses generics to work with any model type (api.User, api.Server, etc.).
func unmarshalList[T any](body io.Reader) ([]*T, *api.Meta, error) {
	// Create an instance of our generic response wrapper.
	// We pass the type T to it.
	response := &api.PaginatedResponse[T]{}

	// Decode the entire JSON response into our struct.
	if err := json.NewDecoder(body).Decode(response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode api list response: %w", err)
	}

	// Flatten the nested structure into a simple slice of models.
	// This is the logic you wanted to avoid repeating!
	results := make([]*T, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}

	// Return the flattened list, the pagination metadata, and no error.
	return results, &response.Meta, nil
}
