package testutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func LoadTestData(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read test data file: "+filePath)
	return data
}
