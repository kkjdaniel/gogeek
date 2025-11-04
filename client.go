package gogeek

import (
	"go.uber.org/ratelimit"
)

// AuthMode represents the authentication method for the client
type AuthMode int

const (
	// AuthNone indicates no authentication
	AuthNone AuthMode = iota
	// AuthAPIKey indicates API key authentication using Bearer token
	AuthAPIKey
	// AuthCookie indicates cookie-based authentication
	AuthCookie
)

// Client represents a GoGeek API client with configurable authentication
type Client struct {
	limiter      ratelimit.Limiter
	authMode     AuthMode
	apiKey       string
	cookieString string
}

// Limiter returns the rate limiter for this client
func (c *Client) Limiter() ratelimit.Limiter {
	return c.limiter
}

// AuthMode returns the authentication mode for this client
func (c *Client) AuthMode() AuthMode {
	return c.authMode
}

// APIKey returns the API key for this client (only valid when AuthMode is AuthAPIKey)
func (c *Client) APIKey() string {
	return c.apiKey
}

// CookieString returns the cookie string for this client (only valid when AuthMode is AuthCookie)
func (c *Client) CookieString() string {
	return c.cookieString
}

// ClientOption is a functional option for configuring a Client
type ClientOption func(*Client)

// WithAPIKey configures the client to use API key authentication
// The API key will be sent as a Bearer token in the Authorization header
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.authMode = AuthAPIKey
		c.apiKey = key
	}
}

// WithCookie configures the client to use cookie-based authentication
// The cookie string should be the raw Cookie header value (e.g., "session=abc; token=xyz")
func WithCookie(cookie string) ClientOption {
	return func(c *Client) {
		c.authMode = AuthCookie
		c.cookieString = cookie
	}
}

// NewClient creates a new GoGeek API client with optional configuration
// By default, the client uses no authentication and has a rate limit of 2 requests per second
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		limiter:  ratelimit.New(2, ratelimit.WithoutSlack),
		authMode: AuthNone,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
