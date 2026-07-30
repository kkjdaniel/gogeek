//go:build !testing

package testutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadTestData reads and returns the test fixture at filePath, failing the test on error.
func LoadTestData(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read test data file: "+filePath)
	return data
}
