package request

import (
	"encoding/xml"
	"net/http"
	"testing"
	"time"

	"github.com/kkjdaniel/gogeek/testutils"
	"github.com/stretchr/testify/require"
)

func TestFetchAndUnmarshal_Success(t *testing.T) {
	defer testutils.ActivateMocks()()

	type TestXML struct {
		XMLName xml.Name `xml:"forum"`
		ID      int      `xml:"id,attr"`
		Title   string   `xml:"title"`
	}

	testURL := "https://example.com/api/test"
	mockDataFileValid := `testdata/valid.xml`
	testutils.SetupMockResponder(t, testURL, mockDataFileValid)

	var result TestXML
	err := FetchAndUnmarshal(testURL, &result)

	require.NoError(t, err, "FetchAndUnmarshal should not return an error with valid XML")
	require.Equal(t, 123, result.ID, "ID should match expected value")
	require.Equal(t, "Example Forum", result.Title, "Title should match expected value")
}

func TestFetchAndUnmarshal_HTTPError(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://nonexistent.example.com"
	testutils.SetupHTTPErrorMock(t, testURL)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when HTTP request fails")
	require.Contains(t, err.Error(), "failed to fetch data from BGG API")
}

func TestFetchAndUnmarshal_BadStatusCode(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/not-found"
	testutils.SetupMockResponderWithStatus(t, testURL, "", http.StatusNotFound)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when status is not 200")
	require.Contains(t, err.Error(), "unexpected status code: 404")
}

func TestFetchAndUnmarshal_InvalidXML(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/invalid-xml"
	invalidXML := "This is not valid XML"
	testutils.SetupMockResponderWithBody(t, testURL, invalidXML, http.StatusOK)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error with invalid XML")
	require.Contains(t, err.Error(), "failed to parse XML")
}

func TestFetchAndUnmarshal_UnmarshalError(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/unmarshal-error"
	validButIncompatibleXML := `<?xml version="1.0"?><different><structure>test</structure></different>`
	testutils.SetupMockResponderWithBody(t, testURL, validButIncompatibleXML, http.StatusOK)

	type MismatchStruct struct {
		XMLName   xml.Name `xml:"expected"`
		SomeField string   `xml:"someField"`
	}

	var result MismatchStruct
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when unmarshaling fails")
	require.Contains(t, err.Error(), "failed to unmarshal XML")
}

func TestFetchAndUnmarshal_Status202_EventualSuccess(t *testing.T) {
	defer testutils.ActivateMocks()()

	originalRetryDelay := retryDelay
	retryDelay = 100 * time.Millisecond
	defer func() { retryDelay = originalRetryDelay }()

	type TestXML struct {
		XMLName xml.Name `xml:"forum"`
		ID      int      `xml:"id,attr"`
		Title   string   `xml:"title"`
	}

	testURL := "https://example.com/api/queued-request"
	mockDataFileValid := `testdata/valid.xml`

	testutils.SetupSequentialResponders(t, testURL, []testutils.MockResponse{
		{StatusCode: http.StatusAccepted, Body: ""},
		{StatusCode: http.StatusAccepted, Body: ""},
		{StatusCode: http.StatusOK, FilePath: mockDataFileValid},
	})

	var result TestXML
	err := FetchAndUnmarshal(testURL, &result)

	require.NoError(t, err, "FetchAndUnmarshal should eventually succeed after 202 responses")
	require.Equal(t, 123, result.ID, "ID should match expected value")
	require.Equal(t, "Example Forum", result.Title, "Title should match expected value")
}

func TestFetchAndUnmarshal_Status202_ExceedsRetries(t *testing.T) {
	defer testutils.ActivateMocks()()

	// Override maxRetries and retryDelay for testing
	originalMaxRetries := maxRetries
	originalRetryDelay := retryDelay
	maxRetries = 3
	retryDelay = 100 * time.Millisecond
	defer func() {
		maxRetries = originalMaxRetries
		retryDelay = originalRetryDelay
	}()

	testURL := "https://example.com/api/always-queued"

	// Setup responders that always return 202
	responses := make([]testutils.MockResponse, maxRetries+1)
	for i := range responses {
		responses[i] = testutils.MockResponse{StatusCode: http.StatusAccepted, Body: ""}
	}
	testutils.SetupSequentialResponders(t, testURL, responses)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should fail after exceeding retries")
	require.Contains(t, err.Error(), "exceeded maximum retries")
}
