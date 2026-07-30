//go:build !testing

package testutils

import (
	"errors"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// MockResponse describes a single mocked HTTP response: its status code,
// body (inline or loaded from FilePath), and any headers to set.
type MockResponse struct {
	StatusCode int
	Body       string
	FilePath   string
	Headers    map[string]string
}

// SetupSequentialResponders registers a mock responder that returns the given
// responses in order, one per call, failing with a 500 once they are exhausted.
func SetupSequentialResponders(t *testing.T, url string, responses []MockResponse) {
	var mu sync.Mutex
	callCount := 0

	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			defer mu.Unlock()

			if callCount >= len(responses) {
				return httpmock.NewStringResponse(500, "No more mock responses configured"), nil
			}

			response := responses[callCount]
			callCount++

			var body string
			if response.FilePath != "" {
				data, err := os.ReadFile(response.FilePath)
				if err != nil {
					t.Fatalf("Failed to read mock data file %s: %v", response.FilePath, err)
				}
				body = string(data)
			} else {
				body = response.Body
			}

			resp := httpmock.NewStringResponse(response.StatusCode, body)

			for key, value := range response.Headers {
				resp.Header.Add(key, value)
			}

			return resp, nil
		},
	)
}

// SetupMockResponder registers a mock responder that returns the contents of
// mockDataPath with a 200 status, and returns the loaded data for assertions.
func SetupMockResponder(t *testing.T, url string, mockDataPath string) []byte {
	mockData, err := os.ReadFile(mockDataPath)
	require.NoError(t, err, "Failed to read mock data file")

	// Register responder
	httpmock.RegisterResponder("GET", url,
		httpmock.NewBytesResponder(200, mockData))

	return mockData
}

// SetupHTTPErrorMock registers a mock responder that fails with a transport-level error.
func SetupHTTPErrorMock(t *testing.T, url string) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("failed to fetch data from BGG API")
		})
}

// SetupMockResponderWithStatus registers a mock responder that returns data with the given status code.
func SetupMockResponderWithStatus(t *testing.T, url string, data string, statusCode int) {
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(statusCode, data))
}

// SetupMockResponderWithBody registers a mock responder that returns body with the given status code.
func SetupMockResponderWithBody(t *testing.T, url string, body string, statusCode int) {
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(statusCode, body))
}

// SetupCountingResponder responds with the given status/body and increments
// *count on every call, so tests can assert how many attempts were made.
func SetupCountingResponder(t *testing.T, url string, statusCode int, body string, count *int) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			*count++
			return httpmock.NewStringResponse(statusCode, body), nil
		})
}

// SetupHeaderCaptureResponder responds 200 with body and passes each request
// to inspect, so tests can assert on request headers.
func SetupHeaderCaptureResponder(t *testing.T, url string, body string, inspect func(*http.Request)) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			inspect(req)
			return httpmock.NewStringResponse(200, body), nil
		})
}

// TestRequestError runs a subtest asserting that queryFunc returns an error
// and a nil result when the underlying HTTP request fails.
func TestRequestError[T any](t *testing.T, url string, queryFunc func(string) (*T, error)) {
	t.Run("Handles errors", func(t *testing.T) {
		defer ActivateMocks()()

		SetupHTTPErrorMock(t, url)

		result, err := queryFunc(url)

		require.Error(t, err, "Function should return an error when request fails")
		require.Nil(t, result, "Result should be nil when an error occurs")
	})
}

// ActivateMocks activates httpmock and returns a cleanup function that
// deactivates it and resets all registered responders.
func ActivateMocks() func() {
	httpmock.Activate()
	return httpmock.DeactivateAndReset
}
