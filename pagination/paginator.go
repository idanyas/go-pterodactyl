// Package pagination provides a generic paginator for the Pterodactyl API.
package pagination

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/idanyas/go-pterodactyl/models"
)

// Paginator provides an iterator-style interface for paginated API results.
type Paginator[T any] struct {
	client      PaginatorClient
	path        string
	currentPage int
	totalPages  int
	perPage     int
	options     ListOptions
}

// PaginatorClient defines the interface required by the Paginator to make API calls.
// This is typically implemented by the main pterodactyl.Client.
type PaginatorClient interface {
	Do(ctx context.Context, method, path string, body, v interface{}) (*http.Response, error)
}

// ListOptions specifies the optional parameters to list methods.
type ListOptions struct {
	Page    int               // Page number to retrieve.
	PerPage int               // Number of items to retrieve per page (max 100).
	Include []string          // Sub-resources to include in the response.
	Filter  map[string]string // Filters to apply to the query.
}

// toQuery converts ListOptions to URL query values.
func (o *ListOptions) toQuery() url.Values {
	v := url.Values{}
	if o.Page > 0 {
		v.Set("page", strconv.Itoa(o.Page))
	}
	if o.PerPage > 0 {
		v.Set("per_page", strconv.Itoa(o.PerPage))
	}
	if len(o.Include) > 0 {
		v.Set("include", strings.Join(o.Include, ","))
	}
	for key, val := range o.Filter {
		v.Set(fmt.Sprintf("filter[%s]", key), val)
	}
	return v
}

// Response is the generic structure for a paginated API response.
type Response[T any] struct {
	Object string      `json:"object"`
	Data   []T         `json:"data"`
	Meta   models.Meta `json:"meta"`
}

// New creates a new Paginator. It performs an initial request to fetch the first page
// and populate the pagination metadata.
func New[T any](ctx context.Context, client PaginatorClient, path string, options ListOptions) ([]T, *Paginator[T], error) {
	// Validate options
	if options.PerPage < 0 {
		return nil, nil, fmt.Errorf("per_page must be non-negative, got %d", options.PerPage)
	}
	if options.PerPage > 100 {
		return nil, nil, fmt.Errorf("per_page must not exceed 100, got %d", options.PerPage)
	}
	if options.Page < 0 {
		return nil, nil, fmt.Errorf("page must be non-negative, got %d", options.Page)
	}

	if options.PerPage == 0 {
		options.PerPage = 50
	}
	if options.Page == 0 {
		options.Page = 1
	}

	p := &Paginator[T]{
		client:      client,
		path:        path,
		currentPage: options.Page,
		perPage:     options.PerPage,
		options:     options,
	}

	items, meta, err := p.fetchPage(ctx, options.Page)
	if err != nil {
		return nil, nil, err
	}

	p.totalPages = meta.Pagination.TotalPages

	// Handle edge case where totalPages is 0 but we got results
	if p.totalPages == 0 && len(items) > 0 {
		p.totalPages = 1
	}

	return items, p, nil
}

// HasMorePages returns true if there are more pages of results to retrieve.
func (p *Paginator[T]) HasMorePages() bool {
	return p.currentPage < p.totalPages
}

// NextPage fetches the next page of results. It returns the items from the next page
// and an error if the request fails. If there are no more pages, it returns nil, nil.
func (p *Paginator[T]) NextPage(ctx context.Context) ([]T, error) {
	if !p.HasMorePages() {
		return nil, nil
	}
	nextPage := p.currentPage + 1
	items, meta, err := p.fetchPage(ctx, nextPage)
	if err != nil {
		return nil, err
	}
	p.currentPage = nextPage

	// Update totalPages in case it changed
	if meta.Pagination.TotalPages > 0 {
		p.totalPages = meta.Pagination.TotalPages
	}

	return items, nil
}

// CurrentPage returns the current page number.
func (p *Paginator[T]) CurrentPage() int {
	return p.currentPage
}

// TotalPages returns the total number of pages.
func (p *Paginator[T]) TotalPages() int {
	return p.totalPages
}

// PerPage returns the number of items per page.
func (p *Paginator[T]) PerPage() int {
	return p.perPage
}

func (p *Paginator[T]) fetchPage(ctx context.Context, page int) ([]T, models.Meta, error) {
	if page <= 0 {
		return nil, models.Meta{}, fmt.Errorf("page must be positive, got %d", page)
	}

	query := p.options.toQuery()
	query.Set("page", strconv.Itoa(page))
	query.Set("per_page", strconv.Itoa(p.perPage))

	fullPath := fmt.Sprintf("%s?%s", p.path, query.Encode())

	var resp Response[struct {
		Object     string `json:"object"`
		Attributes T      `json:"attributes"`
	}]

	_, err := p.client.Do(ctx, http.MethodGet, fullPath, nil, &resp)
	if err != nil {
		return nil, models.Meta{}, fmt.Errorf("failed to fetch page %d: %w", page, err)
	}

	items := make([]T, len(resp.Data))
	for i, item := range resp.Data {
		items[i] = item.Attributes
	}

	return items, resp.Meta, nil
}
