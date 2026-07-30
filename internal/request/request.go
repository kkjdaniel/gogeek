// Package request implements the shared HTTP fetch, retry, and XML
// unmarshalling layer used by every GoGeek endpoint package.
package request

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	gogeek "github.com/kkjdaniel/gogeek/v3"
)

// maxBackoff caps the delay between retries, including server-supplied
// Retry-After values.
const maxBackoff = 30 * time.Second

var (
	cdataRegex       = regexp.MustCompile(`(?s)<!\[CDATA\[(.*?)\]\]>`)
	entityRegex      = regexp.MustCompile(`&(amp|lt|gt|apos|quot|#[0-9]+|#x[0-9a-fA-F]+);`)
	htmlEntityRegex  = regexp.MustCompile(`&([a-zA-Z]+);`)
	controlCharRegex = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F]`)
)

var (
	// ErrEmptyResponse is returned when the response body is empty
	ErrEmptyResponse = errors.New("empty response body")
	// ErrHTTPError is returned when the HTTP request fails
	ErrHTTPError = errors.New("HTTP request failed")
	// ErrUnexpectedStatusCode is returned when the HTTP status code is not 200
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	// ErrMaxRetriesExceeded is returned when the maximum number of retries is exceeded
	ErrMaxRetriesExceeded = errors.New("exceeded maximum retries while waiting for BGG to process request")
	// ErrUnmarshalError is returned when the XML response cannot be unmarshalled
	ErrUnmarshalError = errors.New("failed to unmarshal XML response")
)

// FetchAndUnmarshal performs an authenticated, rate-limited GET request to
// url and unmarshals the XML response into v. It retries on 202 (BGG still
// processing), 429, and 503, honouring Retry-After when present and
// otherwise applying exponential backoff with jitter. All other non-200
// statuses and transport errors fail immediately.
func FetchAndUnmarshal(ctx context.Context, client *gogeek.Client, url string, v any) error {
	maxRetries, baseDelay := client.RetryPolicy()

	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrHTTPError, err)
		}

		if err := client.Prepare(ctx, req); err != nil {
			return err
		}

		resp, err := client.HTTPClient().Do(req)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrHTTPError, err)
		}

		if isRetryable(resp.StatusCode) {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			_ = resp.Body.Close()

			if attempt == maxRetries {
				// 202 means BGG accepted the request and is still preparing
				// it — polling again later may succeed. 429/503 mean the
				// server refused throughout.
				if resp.StatusCode == http.StatusAccepted {
					return ErrMaxRetriesExceeded
				}
				return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
			}

			if err := sleepCtx(ctx, backoffDelay(baseDelay, attempt, retryAfter)); err != nil {
				return err
			}
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
		}

		if readErr != nil {
			return fmt.Errorf("%w: failed to read response body: %v", ErrHTTPError, readErr)
		}

		if len(body) == 0 {
			return ErrEmptyResponse
		}

		return unmarshalXML(body, v)
	}
}

func isRetryable(statusCode int) bool {
	switch statusCode {
	case http.StatusAccepted, http.StatusTooManyRequests, http.StatusServiceUnavailable:
		return true
	}
	return false
}

// backoffDelay computes the wait before the next attempt: the server's
// Retry-After when given, otherwise exponential backoff with equal jitter,
// both capped at maxBackoff.
func backoffDelay(base time.Duration, attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return min(retryAfter, maxBackoff)
	}

	delay := base << attempt
	if delay <= 0 || delay > maxBackoff {
		delay = maxBackoff
	}
	half := delay / 2
	return half + rand.N(half+1)
}

// parseRetryAfter handles both forms of the header: delay-seconds and HTTP-date.
func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}
	if secs, err := strconv.Atoi(header); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(header); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func unmarshalXML(body []byte, v any) error {
	body = fixMalformedXML(body)

	if err := xml.Unmarshal(body, v); err != nil {
		return fmt.Errorf("%w: failed to unmarshal into %T: %v", ErrUnmarshalError, v, err)
	}

	return nil
}

func fixMalformedXML(data []byte) []byte {
	xmlStr := string(data)

	cdataSections := make(map[string]string)
	xmlStr = cdataRegex.ReplaceAllStringFunc(xmlStr, func(match string) string {
		placeholder := fmt.Sprintf("CDATA_PLACEHOLDER_%d", len(cdataSections))
		cdataSections[placeholder] = match
		return placeholder
	})

	xmlStr = entityRegex.ReplaceAllStringFunc(xmlStr, func(s string) string {
		return "ENTITY_PLACEHOLDER" + s[1:]
	})

	xmlStr = htmlEntityRegex.ReplaceAllStringFunc(xmlStr, func(s string) string {
		unescaped := html.UnescapeString(s)
		if unescaped != s {
			return unescaped
		}
		return s
	})

	xmlStr = strings.ReplaceAll(xmlStr, "&", "&amp;")

	xmlStr = strings.ReplaceAll(xmlStr, "ENTITY_PLACEHOLDER", "&")

	for placeholder, cdata := range cdataSections {
		xmlStr = strings.ReplaceAll(xmlStr, placeholder, cdata)
	}

	xmlStr = controlCharRegex.ReplaceAllString(xmlStr, "")

	return []byte(xmlStr)
}
