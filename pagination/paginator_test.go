package pagination

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type mockClient struct {
	handler http.HandlerFunc
}

func (m *mockClient) Do(ctx context.Context, method, path string, body, v interface{}) (*http.Response, error) {
	req := httptest.NewRequest(method, "http://localhost/"+path, nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	m.handler(w, req)
	resp := w.Result()

	if v != nil {
		defer resp.Body.Close()
		// ignore EOF errors caused by empty response bodies
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil && err != io.EOF {
			return resp, err
		}
	}

	return resp, nil
}

type testItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestPaginator(t *testing.T) {
	var requestCount int
	handler := func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, _ := strconv.Atoi(pageStr)

		if page > 2 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"object":"list","data":[],"meta":{"pagination":{"total":20,"count":0,"per_page":10,"current_page":3,"total_pages":2,"links":{}}}}`)
			return
		}

		var items []string
		for i := 1; i <= 10; i++ {
			id := (page-1)*10 + i
			items = append(items, fmt.Sprintf(`{"object":"item","attributes":{"id":%d,"name":"item-%d"}}`, id, id))
		}

		response := fmt.Sprintf(`
		{
			"object": "list",
			"data": [%s],
			"meta": {
				"pagination": {
					"total": 20,
					"count": 10,
					"per_page": 10,
					"current_page": %d,
					"total_pages": 2
				}
			}
		}`, strings.Join(items, ","), page)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	}

	client := &mockClient{handler: handler}
	ctx := context.Background()
	options := ListOptions{PerPage: 10}

	// Initial fetch
	initialItems, paginator, err := New[testItem](ctx, client, "test", options)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("expected 1 request, got %d", requestCount)
	}
	if len(initialItems) != 10 {
		t.Errorf("expected 10 items, got %d", len(initialItems))
	}
	if initialItems[0].ID != 1 || initialItems[9].ID != 10 {
		t.Errorf("unexpected item IDs on first page: got %d and %d", initialItems[0].ID, initialItems[9].ID)
	}
	if !paginator.HasMorePages() {
		t.Error("expected HasMorePages() to be true")
	}

	// Fetch next page
	nextItems, err := paginator.NextPage(ctx)
	if err != nil {
		t.Fatalf("NextPage() failed: %v", err)
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests, got %d", requestCount)
	}
	if len(nextItems) != 10 {
		t.Errorf("expected 10 items on next page, got %d", len(nextItems))
	}
	if nextItems[0].ID != 11 || nextItems[9].ID != 20 {
		t.Errorf("unexpected item IDs on second page: got %d and %d", nextItems[0].ID, nextItems[9].ID)
	}
	if paginator.HasMorePages() {
		t.Error("expected HasMorePages() to be false after fetching last page")
	}

	// Try to fetch past the end
	finalItems, err := paginator.NextPage(ctx)
	if err != nil {
		t.Fatalf("NextPage() beyond end failed: %v", err)
	}
	if requestCount != 2 { // Should not make a new request
		t.Errorf("expected 2 requests, got %d", requestCount)
	}
	if finalItems != nil {
		t.Errorf("expected nil items when fetching beyond end, got %v", finalItems)
	}
}
