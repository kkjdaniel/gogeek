package gogeek

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient_Default(t *testing.T) {
	client := NewClient()

	require.NotNil(t, client, "Client should not be nil")
	require.NotNil(t, client.Limiter(), "Rate limiter should not be nil")
	require.Equal(t, AuthNone, client.AuthMode(), "Default auth mode should be AuthNone")
	require.Equal(t, "", client.APIKey(), "Default API key should be empty")
	require.Equal(t, "", client.CookieString(), "Default cookie string should be empty")
}

func TestNewClient_WithAPIKey(t *testing.T) {
	apiKey := "test-api-key-123"
	client := NewClient(WithAPIKey(apiKey))

	require.NotNil(t, client, "Client should not be nil")
	require.Equal(t, AuthAPIKey, client.AuthMode(), "Auth mode should be AuthAPIKey")
	require.Equal(t, apiKey, client.APIKey(), "API key should match")
	require.Equal(t, "", client.CookieString(), "Cookie string should be empty")
}

func TestNewClient_WithCookie(t *testing.T) {
	cookie := "bggusername=test; bggpassword=abc123; SessionID=xyz789"
	client := NewClient(WithCookie(cookie))

	require.NotNil(t, client, "Client should not be nil")
	require.Equal(t, AuthCookie, client.AuthMode(), "Auth mode should be AuthCookie")
	require.Equal(t, cookie, client.CookieString(), "Cookie string should match")
	require.Equal(t, "", client.APIKey(), "API key should be empty")
}

func TestNewClient_RateLimiter(t *testing.T) {
	client := NewClient()

	// Rate limiter should allow at least one call immediately
	start := time.Now()
	client.Limiter().Take()
	duration := time.Since(start)

	// First call should be nearly instant
	require.Less(t, duration, 100*time.Millisecond, "First rate limit call should be instant")

	// Second call should be delayed by approximately 0.5 seconds (2 requests per second)
	start = time.Now()
	client.Limiter().Take()
	duration = time.Since(start)

	// Should wait close to 0.5 seconds (2 requests per second = 0.5s between requests)
	require.Greater(t, duration, 450*time.Millisecond, "Second call should wait ~0.5 seconds")
	require.Less(t, duration, 600*time.Millisecond, "Second call should not wait too long")
}
