package request

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"testing"
	"testing/iotest"
	"time"

	"github.com/jarcoal/httpmock"
	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/internal/testutils"
	"github.com/stretchr/testify/require"
)

// testClient returns a client tuned for fast tests: effectively no rate
// limiting and short retry delays.
func testClient(opts ...gogeek.ClientOption) *gogeek.Client {
	base := []gogeek.ClientOption{
		gogeek.WithRateLimit(1000),
		gogeek.WithRetry(5, 10*time.Millisecond),
	}
	return gogeek.NewClient(gogeek.APIKey("test-key"), append(base, opts...)...)
}

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
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.NoError(t, err, "FetchAndUnmarshal should not return an error with valid XML")
	require.Equal(t, 123, result.ID, "ID should match expected value")
	require.Equal(t, "Example Forum", result.Title, "Title should match expected value")
}

func TestFetchAndUnmarshal_SendsAuthHeader(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/auth-check"
	var gotAuth string
	testutils.SetupHeaderCaptureResponder(t, testURL, `<item id="1"></item>`, func(req *http.Request) {
		gotAuth = req.Header.Get("Authorization")
	})

	var result struct {
		ID int `xml:"id,attr"`
	}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.NoError(t, err)
	require.Equal(t, "Bearer test-key", gotAuth, "Request should carry the API key as a Bearer token")
}

func TestFetchAndUnmarshal_HTTPError(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://nonexistent.example.com"
	testutils.SetupHTTPErrorMock(t, testURL)

	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when HTTP request fails")
	require.True(t, errors.Is(err, ErrHTTPError), "Error should be of type ErrHTTPError")
}

func TestFetchAndUnmarshal_BadStatusCode(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/not-found"
	testutils.SetupMockResponderWithStatus(t, testURL, "", http.StatusNotFound)

	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when status is not 200")
	require.True(t, errors.Is(err, ErrUnexpectedStatusCode), "Error should be of type ErrUnexpectedStatusCode")
	require.Contains(t, err.Error(), "404", "Error should include the status code")
}

func TestFetchAndUnmarshal_ServerErrorFailsFast(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/server-error"
	callCount := 0
	testutils.SetupCountingResponder(t, testURL, http.StatusInternalServerError, "", &callCount)

	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrUnexpectedStatusCode), "Error should be of type ErrUnexpectedStatusCode")
	require.Equal(t, 1, callCount, "500 responses should not be retried")
}

func TestFetchAndUnmarshal_InvalidXML(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/invalid-xml"
	invalidXML := "This is not valid XML"
	testutils.SetupMockResponderWithBody(t, testURL, invalidXML, http.StatusOK)

	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error with invalid XML")
	require.True(t, errors.Is(err, ErrUnmarshalError), "Error should be of type ErrUnmarshalError")
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
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when unmarshaling fails")
	require.True(t, errors.Is(err, ErrUnmarshalError), "Error should be of type ErrUnmarshalError")
}

func TestFetchAndUnmarshal_InvalidURL(t *testing.T) {
	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), "http://example.com/\x7f", &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when the request cannot be built")
	require.True(t, errors.Is(err, ErrHTTPError), "Error should be of type ErrHTTPError")
}

func TestFetchAndUnmarshal_EmptyBody(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/empty"
	testutils.SetupMockResponderWithBody(t, testURL, "", http.StatusOK)

	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.True(t, errors.Is(err, ErrEmptyResponse), "Error should be of type ErrEmptyResponse")
}

func TestFetchAndUnmarshal_BodyReadError(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/read-error"
	httpmock.RegisterResponder("GET", testURL,
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(iotest.ErrReader(errors.New("connection reset"))),
				Header:     http.Header{},
				Request:    req,
			}, nil
		})

	var result struct{}
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should return an error when the body cannot be read")
	require.True(t, errors.Is(err, ErrHTTPError), "Error should be of type ErrHTTPError")
	require.Contains(t, err.Error(), "connection reset", "Error should include the underlying read error")
}

func TestFetchAndUnmarshal_Status202_EventualSuccess(t *testing.T) {
	defer testutils.ActivateMocks()()

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
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.NoError(t, err, "FetchAndUnmarshal should eventually succeed after 202 responses")
	require.Equal(t, 123, result.ID, "ID should match expected value")
	require.Equal(t, "Example Forum", result.Title, "Title should match expected value")
}

func TestFetchAndUnmarshal_Status202_ExceedsRetries(t *testing.T) {
	defer testutils.ActivateMocks()()

	const maxRetries = 3
	testURL := "https://example.com/api/always-queued"

	responses := make([]testutils.MockResponse, maxRetries+1)
	for i := range responses {
		responses[i] = testutils.MockResponse{StatusCode: http.StatusAccepted, Body: ""}
	}
	testutils.SetupSequentialResponders(t, testURL, responses)

	var result struct{}
	client := testClient(gogeek.WithRetry(maxRetries, 10*time.Millisecond))
	err := FetchAndUnmarshal(context.Background(), client, testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should fail after exceeding retries")
	require.True(t, errors.Is(err, ErrMaxRetriesExceeded), "Error should be of type ErrMaxRetriesExceeded")
	require.Contains(t, err.Error(), "exceeded maximum retries")
}

func TestFetchAndUnmarshal_Status429_EventualSuccess(t *testing.T) {
	defer testutils.ActivateMocks()()

	type TestXML struct {
		XMLName xml.Name `xml:"forum"`
		ID      int      `xml:"id,attr"`
	}

	testURL := "https://example.com/api/throttled"
	testutils.SetupSequentialResponders(t, testURL, []testutils.MockResponse{
		{StatusCode: http.StatusTooManyRequests, Body: ""},
		{StatusCode: http.StatusOK, FilePath: `testdata/valid.xml`},
	})

	var result TestXML
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.NoError(t, err, "FetchAndUnmarshal should retry 429 and succeed")
	require.Equal(t, 123, result.ID)
}

func TestFetchAndUnmarshal_Status429_ExceedsRetries(t *testing.T) {
	defer testutils.ActivateMocks()()

	const maxRetries = 2
	testURL := "https://example.com/api/always-throttled"

	responses := make([]testutils.MockResponse, maxRetries+1)
	for i := range responses {
		responses[i] = testutils.MockResponse{StatusCode: http.StatusTooManyRequests, Body: ""}
	}
	testutils.SetupSequentialResponders(t, testURL, responses)

	var result struct{}
	client := testClient(gogeek.WithRetry(maxRetries, 10*time.Millisecond))
	err := FetchAndUnmarshal(context.Background(), client, testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should fail after exhausting 429 retries")
	require.True(t, errors.Is(err, ErrUnexpectedStatusCode),
		"429 exhaustion should wrap ErrUnexpectedStatusCode, not ErrMaxRetriesExceeded")
	require.Contains(t, err.Error(), "429", "Error should include the status code")
}

func TestFetchAndUnmarshal_Status503_EventualSuccess(t *testing.T) {
	defer testutils.ActivateMocks()()

	type TestXML struct {
		XMLName xml.Name `xml:"forum"`
		ID      int      `xml:"id,attr"`
	}

	testURL := "https://example.com/api/unavailable"
	testutils.SetupSequentialResponders(t, testURL, []testutils.MockResponse{
		{StatusCode: http.StatusServiceUnavailable, Body: ""},
		{StatusCode: http.StatusOK, FilePath: `testdata/valid.xml`},
	})

	var result TestXML
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)

	require.NoError(t, err, "FetchAndUnmarshal should retry 503 and succeed")
	require.Equal(t, 123, result.ID)
}

func TestFetchAndUnmarshal_RetryAfterHeader(t *testing.T) {
	defer testutils.ActivateMocks()()

	type TestXML struct {
		XMLName xml.Name `xml:"forum"`
		ID      int      `xml:"id,attr"`
	}

	testURL := "https://example.com/api/retry-after"
	testutils.SetupSequentialResponders(t, testURL, []testutils.MockResponse{
		{StatusCode: http.StatusTooManyRequests, Body: "", Headers: map[string]string{"Retry-After": "1"}},
		{StatusCode: http.StatusOK, FilePath: `testdata/valid.xml`},
	})

	var result TestXML
	start := time.Now()
	err := FetchAndUnmarshal(context.Background(), testClient(), testURL, &result)
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.GreaterOrEqual(t, elapsed, 1*time.Second, "Retry-After: 1 should delay the retry by at least a second")
}

func TestFetchAndUnmarshal_CancelledContext(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/cancelled"
	testutils.SetupMockResponderWithBody(t, testURL, "<item/>", http.StatusOK)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result struct{}
	err := FetchAndUnmarshal(ctx, testClient(), testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should fail with a cancelled context")
	require.True(t, errors.Is(err, context.Canceled), "Error should be context.Canceled")
}

func TestFetchAndUnmarshal_ContextCancelledDuringRetry(t *testing.T) {
	defer testutils.ActivateMocks()()

	testURL := "https://example.com/api/cancel-mid-retry"
	testutils.SetupMockResponderWithStatus(t, testURL, "", http.StatusAccepted)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var result struct{}
	client := testClient(gogeek.WithRetry(5, 10*time.Second))
	err := FetchAndUnmarshal(ctx, client, testURL, &result)

	require.Error(t, err, "FetchAndUnmarshal should abort the retry sleep when ctx expires")
	require.True(t, errors.Is(err, context.DeadlineExceeded), "Error should be context.DeadlineExceeded")
}

func TestParseRetryAfter(t *testing.T) {
	require.Equal(t, time.Duration(0), parseRetryAfter(""))
	require.Equal(t, 5*time.Second, parseRetryAfter("5"))
	require.Equal(t, time.Duration(0), parseRetryAfter("-1"))
	require.Equal(t, time.Duration(0), parseRetryAfter("garbage"))

	future := time.Now().Add(10 * time.Second).UTC().Format(http.TimeFormat)
	d := parseRetryAfter(future)
	require.Greater(t, d, 8*time.Second, "HTTP-date Retry-After should yield the remaining duration")
	require.LessOrEqual(t, d, 10*time.Second)

	past := time.Now().Add(-10 * time.Second).UTC().Format(http.TimeFormat)
	require.Equal(t, time.Duration(0), parseRetryAfter(past), "Past HTTP-date should yield zero")
}

func TestBackoffDelay(t *testing.T) {
	base := 2 * time.Second

	// Retry-After wins, capped at maxBackoff.
	require.Equal(t, 5*time.Second, backoffDelay(base, 0, 5*time.Second))
	require.Equal(t, maxBackoff, backoffDelay(base, 0, time.Hour))

	// Exponential growth with jitter: delay for attempt n is in [d/2, d]
	// where d = base * 2^n, capped at maxBackoff.
	for attempt := 0; attempt < 10; attempt++ {
		expected := base << attempt
		if expected <= 0 || expected > maxBackoff {
			expected = maxBackoff
		}
		got := backoffDelay(base, attempt, 0)
		require.GreaterOrEqual(t, got, expected/2, "attempt %d", attempt)
		require.LessOrEqual(t, got, expected, "attempt %d", attempt)
	}
}

func TestFixMalformedXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"preserves XML entities", "<a>Tom &amp; Jerry &lt;3 &#233; &#xE9;</a>", "<a>Tom &amp; Jerry &lt;3 &#233; &#xE9;</a>"},
		{"unescapes HTML entities", "<a>caf&eacute;</a>", "<a>café</a>"},
		{"escapes unknown entities", "<a>&bogus;</a>", "<a>&amp;bogus;</a>"},
		{"escapes bare ampersands", "<a>Dungeons & Dragons</a>", "<a>Dungeons &amp; Dragons</a>"},
		{"strips control characters", "<a>bad\x00\x0Bchars\ttab</a>", "<a>badchars\ttab</a>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, string(fixMalformedXML([]byte(tt.input))))
		})
	}
}
