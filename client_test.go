package gogeek

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient_Defaults(t *testing.T) {
	client := NewClient(APIKey("test-key"))

	require.NotNil(t, client, "Client should not be nil")
	require.NotNil(t, client.HTTPClient(), "HTTP client should not be nil")
	require.Equal(t, 30*time.Second, client.HTTPClient().Timeout, "Default HTTP timeout should be 30s")

	maxRetries, delay := client.RetryPolicy()
	require.Equal(t, 5, maxRetries, "Default max retries should be 5")
	require.Equal(t, 2*time.Second, delay, "Default retry delay should be 2s")
}

func TestPrepare_APIKey(t *testing.T) {
	client := NewClient(APIKey("test-api-key-123"))

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	require.NoError(t, client.Prepare(context.Background(), req))
	require.Equal(t, "Bearer test-api-key-123", req.Header.Get("Authorization"))
	require.Empty(t, req.Header.Get("Cookie"))
}

func TestPrepare_Cookie(t *testing.T) {
	cookie := "bggusername=test; bggpassword=abc123; SessionID=xyz789"
	client := NewClient(Cookie(cookie))

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	require.NoError(t, client.Prepare(context.Background(), req))
	require.Equal(t, cookie, req.Header.Get("Cookie"))
	require.Empty(t, req.Header.Get("Authorization"))
}

func TestPrepare_CancelledContext(t *testing.T) {
	client := NewClient(APIKey("test-key"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	require.Error(t, client.Prepare(ctx, req), "Prepare should fail with a cancelled context")
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	client := NewClient(APIKey("test-key"), WithHTTPClient(custom))

	require.Same(t, custom, client.HTTPClient(), "Custom HTTP client should be used")
}

func TestWithRetry(t *testing.T) {
	client := NewClient(APIKey("test-key"), WithRetry(2, 100*time.Millisecond))

	maxRetries, delay := client.RetryPolicy()
	require.Equal(t, 2, maxRetries)
	require.Equal(t, 100*time.Millisecond, delay)
}

func TestRateLimiter(t *testing.T) {
	client := NewClient(APIKey("test-key"))
	ctx := context.Background()

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	// First call should be nearly instant
	start := time.Now()
	require.NoError(t, client.Prepare(ctx, req))
	require.Less(t, time.Since(start), 100*time.Millisecond, "First rate limit call should be instant")

	// Second call should wait ~0.5 seconds (2 requests per second)
	start = time.Now()
	require.NoError(t, client.Prepare(ctx, req))
	duration := time.Since(start)
	require.Greater(t, duration, 450*time.Millisecond, "Second call should wait ~0.5 seconds")
	require.Less(t, duration, 600*time.Millisecond, "Second call should not wait too long")
}

func TestWithRateLimit(t *testing.T) {
	client := NewClient(APIKey("test-key"), WithRateLimit(10))
	ctx := context.Background()

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	start := time.Now()
	for i := 0; i < 3; i++ {
		require.NoError(t, client.Prepare(ctx, req))
	}
	duration := time.Since(start)
	require.Greater(t, duration, 150*time.Millisecond, "10 rps should space requests ~100ms apart")
	require.Less(t, duration, 400*time.Millisecond, "10 rps should not wait too long")
}
