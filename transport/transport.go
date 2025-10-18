// Package transport provides an http.RoundTripper for the Pterodactyl client.
package transport

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultMaxRetries       = 3
	defaultRetryWaitMin     = 1 * time.Second
	defaultRetryWaitMax     = 5 * time.Second
	defaultRateLimitMaxWait = 5 * time.Minute
	headerAccept            = "Accept"
	headerAuth              = "Authorization"
	headerUA                = "User-Agent"
)

// RateLimitInfo contains the rate limit information from an API response.
type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// ParseRateLimit extracts rate limit information from response headers.
func ParseRateLimit(r *http.Response) RateLimitInfo {
	var rate RateLimitInfo
	if r == nil || r.Header == nil {
		return rate
	}

	if limit := r.Header.Get("X-RateLimit-Limit"); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get("X-RateLimit-Reset"); reset != "" {
		if v, err := strconv.ParseInt(reset, 10, 64); err == nil {
			rate.Reset = time.Unix(v, 0)
		}
	}
	return rate
}

// TransportOption is a functional option for configuring a Transport.
type TransportOption func(*Transport)

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(max int) TransportOption {
	return func(t *Transport) {
		if max > 0 {
			t.maxRetries = max
		}
	}
}

// WithRetryWaitMin sets the minimum wait time between retries.
func WithRetryWaitMin(d time.Duration) TransportOption {
	return func(t *Transport) {
		if d > 0 {
			t.retryWaitMin = d
		}
	}
}

// WithRetryWaitMax sets the maximum wait time between retries.
func WithRetryWaitMax(d time.Duration) TransportOption {
	return func(t *Transport) {
		if d > 0 {
			t.retryWaitMax = d
		}
	}
}

// WithRateLimitMaxWait sets the maximum time to wait for a rate limit reset.
func WithRateLimitMaxWait(d time.Duration) TransportOption {
	return func(t *Transport) {
		if d > 0 {
			t.rateLimitMaxWait = d
		}
	}
}

// Transport is an http.RoundTripper that handles authentication,
// rate limiting, and retries for the Pterodactyl API.
type Transport struct {
	base      http.RoundTripper
	apiKey    string
	accept    string
	userAgent string

	// Retry configuration
	maxRetries       int
	retryWaitMin     time.Duration
	retryWaitMax     time.Duration
	rateLimitMaxWait time.Duration
}

// New creates a new Transport with optional configuration.
func New(base http.RoundTripper, apiKey, apiVersion, userAgent string, opts ...TransportOption) *Transport {
	t := &Transport{
		base:             base,
		apiKey:           apiKey,
		accept:           fmt.Sprintf("Application/vnd.pterodactyl.%s+json", apiVersion),
		userAgent:        userAgent,
		maxRetries:       defaultMaxRetries,
		retryWaitMin:     defaultRetryWaitMin,
		retryWaitMax:     defaultRetryWaitMax,
		rateLimitMaxWait: defaultRateLimitMaxWait,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// RoundTrip executes a single HTTP transaction, adding required headers and handling retries.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(headerAccept, t.accept)
	req.Header.Set(headerAuth, "Bearer "+t.apiKey)
	req.Header.Set(headerUA, t.userAgent)

	var resp *http.Response
	var err error

	for i := 0; i < t.maxRetries; i++ {
		// Clone the request body if it exists
		if req.Body != nil {
			var bodyErr error
			req.Body, bodyErr = req.GetBody()
			if bodyErr != nil {
				return nil, bodyErr
			}
		}

		resp, err = t.base.RoundTrip(req)
		if err != nil {
			// Network-level error, retry
			if t.waitAndRetry(req.Context(), i) {
				continue
			}
			return nil, err
		}

		// Success (2xx) or a non-retriable error (e.g., 4xx, except 429)
		if resp.StatusCode < http.StatusInternalServerError && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// Handle 429 Too Many Requests
		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimit := ParseRateLimit(resp)

			// Calculate wait duration, ensuring it's not negative or too long
			waitDuration := time.Until(rateLimit.Reset)
			if waitDuration <= 0 {
				waitDuration = t.retryWaitMin
			} else if waitDuration > t.rateLimitMaxWait {
				waitDuration = t.rateLimitMaxWait
			}

			// Drain the body to allow connection reuse
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			select {
			case <-time.After(waitDuration):
				continue // Retry the request
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}

		// Handle 5xx Server Errors
		if resp.StatusCode >= http.StatusInternalServerError {
			// Drain the body before retrying
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if t.waitAndRetry(req.Context(), i) {
				continue
			}

			// Context was cancelled or we're out of retries
			if req.Context().Err() != nil {
				return nil, req.Context().Err()
			}
		}

		// If we are here, it means we ran out of retries for 5xx errors
		break
	}

	return resp, err
}

// waitAndRetry calculates the backoff duration, waits, and returns true if a retry should be attempted.
func (t *Transport) waitAndRetry(ctx context.Context, retryCount int) bool {
	if retryCount >= t.maxRetries-1 {
		return false
	}

	// Exponential backoff with jitter
	backoff := float64(t.retryWaitMin) * math.Pow(2, float64(retryCount))
	if backoff > float64(t.retryWaitMax) {
		backoff = float64(t.retryWaitMax)
	}
	backoff *= (1 + rand.Float64()*0.5) // Add up to 50% jitter
	wait := time.Duration(backoff)

	select {
	case <-time.After(wait):
		return true
	case <-ctx.Done():
		return false
	}
}
