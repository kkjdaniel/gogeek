//go:build !testing

package testutils

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

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

func ActivateMocks() func() {
	httpmock.Activate()
	return httpmock.DeactivateAndReset
}
