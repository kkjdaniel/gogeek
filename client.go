// Package gogeek provides a client for the BoardGameGeek XML API2.
package gogeek

import (
	"context"
	"errors"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// ErrInvalidOption is wrapped by errors returned when a query option is
// given an out-of-range or otherwise invalid argument.
var ErrInvalidOption = errors.New("invalid option")

const (
	defaultTimeout    = 30 * time.Second
	defaultRateLimit  = 2
	defaultMaxRetries = 5
	defaultRetryDelay = 2 * time.Second
)

type authMode int

const (
	authAPIKey authMode = iota + 1
	authCookie
)

// Auth holds credentials for the BGG API. Construct one with APIKey or
// Cookie; the zero value sends no authentication headers, which BGG rejects.
type Auth struct {
	mode  authMode
	value string
}

// APIKey returns an Auth that authenticates requests with the given API key,
// sent as a Bearer token in the Authorization header. API keys can be
// requested at https://boardgamegeek.com/applications.
func APIKey(key string) Auth {
	return Auth{mode: authAPIKey, value: key}
}

// Cookie returns an Auth that authenticates requests with the given raw
// Cookie header value (e.g. "bggusername=user; bggpassword=...; SessionID=...").
// Cookie authentication can access private collection data that an API key
// may not.
func Cookie(cookie string) Auth {
	return Auth{mode: authCookie, value: cookie}
}

// Client is a BGG API client. Create one with NewClient.
type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	auth       Auth
	maxRetries int
	retryDelay time.Duration
}

// ClientOption is a functional option for configuring a Client.
type ClientOption func(*Client)

// WithHTTPClient replaces the default HTTP client (30 second timeout) with a
// custom one, for callers who need a specific transport, proxy, or timeout.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithRateLimit replaces the default rate limit of 2 requests per second.
func WithRateLimit(rps int) ClientOption {
	return func(c *Client) {
		c.limiter = rate.NewLimiter(rate.Limit(rps), 1)
	}
}

// WithRetry replaces the default retry behaviour (5 retries, 2 second base
// delay) used when BGG responds with a retryable status (202, 429, 503).
// The delay is the base for exponential backoff and is superseded by a
// Retry-After response header when present.
func WithRetry(maxRetries int, delay time.Duration) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryDelay = delay
	}
}

// NewClient creates a new BGG API client. All BGG endpoints require
// authentication, so an Auth (from APIKey or Cookie) is mandatory. By
// default the client sends at most 2 requests per second, times out
// requests after 30 seconds, and retries retryable responses up to 5 times;
// use the options to change any of these.
func NewClient(auth Auth, opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		limiter:    rate.NewLimiter(rate.Limit(defaultRateLimit), 1),
		auth:       auth,
		maxRetries: defaultMaxRetries,
		retryDelay: defaultRetryDelay,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Prepare applies authentication headers and rate limiting to req, blocking
// until the rate limiter permits the request or ctx is done. It allows raw
// requests to BGG without exposing the client's credentials.
func (c *Client) Prepare(ctx context.Context, req *http.Request) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	switch c.auth.mode {
	case authAPIKey:
		req.Header.Set("Authorization", "Bearer "+c.auth.value)
	case authCookie:
		req.Header.Set("Cookie", c.auth.value)
	}

	return nil
}

// HTTPClient returns the HTTP client used to send requests, either the
// default or the one supplied via WithHTTPClient.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// RetryPolicy returns the maximum number of retries and the base backoff
// delay used for retryable responses, as configured by WithRetry.
func (c *Client) RetryPolicy() (maxRetries int, delay time.Duration) {
	return c.maxRetries, c.retryDelay
}
