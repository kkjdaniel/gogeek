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

type MockResponse struct {
	StatusCode int
	Body       string
	FilePath   string
	Headers    map[string]string
}

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

func SetupMockResponder(t *testing.T, url string, mockDataPath string) []byte {
	mockData, err := os.ReadFile(mockDataPath)
	require.NoError(t, err, "Failed to read mock data file")

	// Register responder
	httpmock.RegisterResponder("GET", url,
		httpmock.NewBytesResponder(200, mockData))

	return mockData
}

func SetupHTTPErrorMock(t *testing.T, url string) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("failed to fetch data from BGG API")
		})
}

func SetupMockResponderWithStatus(t *testing.T, url string, data string, statusCode int) {
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(statusCode, data))
}

func SetupMockResponderWithBody(t *testing.T, url string, body string, statusCode int) {
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(statusCode, body))
}

func TestRequestError[T any](t *testing.T, url string, queryFunc func(string) (*T, error)) {
	t.Run("Handles errors", func(t *testing.T) {
		defer ActivateMocks()()

		SetupHTTPErrorMock(t, url)

		result, err := queryFunc(url)

		require.Error(t, err, "Function should return an error when request fails")
		require.Nil(t, result, "Result should be nil when an error occurs")
	})
}

func ActivateMocks() func() {
	httpmock.Activate()
	return httpmock.DeactivateAndReset
}
