package testutils

import (
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

func ActivateMocks() func() {
	httpmock.Activate()
	return httpmock.DeactivateAndReset
}
