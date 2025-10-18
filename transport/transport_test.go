package transport

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestTransport_Headers(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get(headerAccept); got != "Application/vnd.pterodactyl.v1+json" {
			t.Errorf("Accept header = %q, want %q", got, "Application/vnd.pterodactyl.v1+json")
		}
		if got := r.Header.Get(headerAuth); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer test-key")
		}
		if got := r.Header.Get(headerUA); got == "" {
			t.Error("User-Agent header is empty")
		}
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tp := New(http.DefaultTransport, "test-key", "v1", "test-agent")
	client := &http.Client{Transport: tp}
	req, _ := http.NewRequest("GET", server.URL, nil)

	_, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do failed: %v", err)
	}
}

func TestTransport_Retry(t *testing.T) {
	var requests int32
	handler := func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requests, 1)
		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tp := New(http.DefaultTransport, "test-key", "v1", "test-agent")
	client := &http.Client{Transport: tp}
	req, _ := http.NewRequest("GET", server.URL, nil)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("client.Do failed: %v", err)
	}
	defer resp.Body.Close()

	if requests != 3 {
		t.Errorf("expected 3 requests, got %d", requests)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %s", resp.Status)
	}
	if duration < 1*time.Second {
		t.Errorf("retry backoff seems too short: %v", duration)
	}
}

func TestTransport_RateLimit(t *testing.T) {
	var requests int32
	resetTime := time.Now().Add(200 * time.Millisecond)

	handler := func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requests, 1)
		if count == 1 {
			w.Header().Set("X-RateLimit-Limit", "240")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tp := New(http.DefaultTransport, "test-key", "v1", "test-agent")
	client := &http.Client{Transport: tp}
	req, _ := http.NewRequest("GET", server.URL, nil)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("client.Do failed: %v", err)
	}
	defer resp.Body.Close()

	if requests != 2 {
		t.Errorf("expected 2 requests, got %d", requests)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %s", resp.Status)
	}
	if duration < 150*time.Millisecond {
		t.Errorf("rate limit wait seems too short: %v (expected at least 150ms)", duration)
	}
}

func TestTransport_Retry_ContextCancel(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tp := New(http.DefaultTransport, "test-key", "v1", "test-agent")
	client := &http.Client{Transport: tp}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)

	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestParseRateLimit(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     RateLimitInfo
	}{
		{
			name: "valid headers",
			response: &http.Response{
				Header: http.Header{
					"X-Ratelimit-Limit":     []string{"240"},
					"X-Ratelimit-Remaining": []string{"100"},
					"X-Ratelimit-Reset":     []string{"1640000000"},
				},
			},
			want: RateLimitInfo{
				Limit:     240,
				Remaining: 100,
				Reset:     time.Unix(1640000000, 0),
			},
		},
		{
			name:     "nil response",
			response: nil,
			want:     RateLimitInfo{},
		},
		{
			name: "nil headers",
			response: &http.Response{
				Header: nil,
			},
			want: RateLimitInfo{},
		},
		{
			name: "empty headers",
			response: &http.Response{
				Header: http.Header{},
			},
			want: RateLimitInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRateLimit(tt.response)

			if got.Limit != tt.want.Limit {
				t.Errorf("Limit = %d, want %d", got.Limit, tt.want.Limit)
			}
			if got.Remaining != tt.want.Remaining {
				t.Errorf("Remaining = %d, want %d", got.Remaining, tt.want.Remaining)
			}
			if !got.Reset.Equal(tt.want.Reset) {
				t.Errorf("Reset = %v (unix: %d), want %v (unix: %d)",
					got.Reset, got.Reset.Unix(),
					tt.want.Reset, tt.want.Reset.Unix())
			}
		})
	}
}

func TestTransport_ConfigurableRetries(t *testing.T) {
	var requests int32
	handler := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Test with custom max retries
	tp := New(
		http.DefaultTransport,
		"test-key",
		"v1",
		"test-agent",
		WithMaxRetries(5),
		WithRetryWaitMin(10*time.Millisecond),
		WithRetryWaitMax(50*time.Millisecond),
	)
	client := &http.Client{Transport: tp}
	req, _ := http.NewRequest("GET", server.URL, nil)

	_, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do failed: %v", err)
	}

	// Should make 5 attempts (maxRetries)
	if requests != 5 {
		t.Errorf("expected 5 requests with custom maxRetries, got %d", requests)
	}
}

func TestTransport_RateLimitMaxWait(t *testing.T) {
	var requests int32
	// Set reset time far in the future (1 hour)
	resetTime := time.Now().Add(1 * time.Hour)

	handler := func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requests, 1)
		if count == 1 {
			w.Header().Set("X-RateLimit-Limit", "240")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Set max wait to 100ms instead of default 5 minutes
	tp := New(
		http.DefaultTransport,
		"test-key",
		"v1",
		"test-agent",
		WithRateLimitMaxWait(100*time.Millisecond),
	)
	client := &http.Client{Transport: tp}
	req, _ := http.NewRequest("GET", server.URL, nil)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("client.Do failed: %v", err)
	}
	defer resp.Body.Close()

	// Should wait for maxWait (100ms) instead of 1 hour
	if duration < 80*time.Millisecond || duration > 200*time.Millisecond {
		t.Errorf("expected wait time around 100ms, got %v", duration)
	}
}
