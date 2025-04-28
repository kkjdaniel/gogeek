package request

import (
	"encoding/xml"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/kkjdaniel/gogeek/testutils"
	"github.com/stretchr/testify/require"
)

func TestFetchAndUnmarshal_Success(t *testing.T) {
	defer testutils.ActivateMocks()()

	type TestXML struct {
		ID    int    `xml:"id,attr"`
		Title string `xml:"title"`
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
	require.True(t, errors.Is(err, ErrHTTPError), "Error should be of type ErrHTTPError")
}

func TestFetchAndUnmarshal_BadStatusCode(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/not-found"
	testutils.SetupMockResponderWithStatus(t, testURL, "", http.StatusNotFound)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when status is not 200")
	require.True(t, errors.Is(err, ErrUnexpectedStatusCode), "Error should be of type ErrUnexpectedStatusCode")
	require.Contains(t, err.Error(), "404", "Error should include the status code")
}

func TestFetchAndUnmarshal_InvalidXML(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/invalid-xml"
	invalidXML := "This is not valid XML"
	testutils.SetupMockResponderWithBody(t, testURL, invalidXML, http.StatusOK)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error with invalid XML")
	require.True(t, errors.Is(err, ErrXMLParseError), "Error should be of type ErrXMLParseError")
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
	require.True(t, errors.Is(err, ErrUnmarshalError), "Error should be of type ErrUnmarshalError")
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

	originalMaxRetries := maxRetries
	originalRetryDelay := retryDelay
	maxRetries = 3
	retryDelay = 100 * time.Millisecond
	defer func() {
		maxRetries = originalMaxRetries
		retryDelay = originalRetryDelay
	}()

	testURL := "https://example.com/api/always-queued"

	responses := make([]testutils.MockResponse, maxRetries+1)
	for i := range responses {
		responses[i] = testutils.MockResponse{StatusCode: http.StatusAccepted, Body: ""}
	}
	testutils.SetupSequentialResponders(t, testURL, responses)

	var result struct{}
	err := FetchAndUnmarshal(testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should fail after exceeding retries")
	require.True(t, errors.Is(err, ErrMaxRetriesExceeded), "Error should be of type ErrMaxRetriesExceeded")
	require.Contains(t, err.Error(), "exceeded maximum retries")
}

func TestFixMalformedXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple unescaped ampersand",
			input:    `<item>Dungeons & Dragons</item>`,
			expected: `<item>Dungeons &amp; Dragons</item>`,
		},
		{
			name:     "Multiple unescaped ampersands",
			input:    `<game>Ticket to Ride: Rails & Sails & More</game>`,
			expected: `<game>Ticket to Ride: Rails &amp; Sails &amp; More</game>`,
		},
		{
			name:     "Preserve existing entities",
			input:    `<description>This game uses &lt;cards&gt; &amp; dice</description>`,
			expected: `<description>This game uses &lt;cards&gt; &amp; dice</description>`,
		},
		{
			name:     "Preserve numeric entities",
			input:    `<text>Copyright &#169; 2025 &amp; trademark &#8482;</text>`,
			expected: `<text>Copyright &#169; 2025 &amp; trademark &#8482;</text>`,
		},
		{
			name:     "Mixed valid and invalid ampersands",
			input:    `<item>&lt;Dungeons & Dragons&gt; uses dice &amp; cards</item>`,
			expected: `<item>&lt;Dungeons &amp; Dragons&gt; uses dice &amp; cards</item>`,
		},
		{
			name:     "Remove control characters",
			input:    "<description>Game\x0Bdescription\x1Fhere</description>",
			expected: "<description>Gamedescriptionhere</description>",
		},
		{
			name:     "Complex mixed case",
			input:    "<item>\x0BDungeons & Dragons &#169; &amp; More\x1F</item>",
			expected: "<item>Dungeons &amp; Dragons &#169; &amp; More</item>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(fixMalformedXML([]byte(tt.input)))
			require.Equal(t, tt.expected, result, "XML should be fixed correctly")

			var anyXML interface{}
			err := xml.Unmarshal([]byte(result), &anyXML)
			require.NoError(t, err, "Fixed XML should be valid")
		})
	}
}
