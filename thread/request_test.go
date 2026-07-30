package thread

import (
	"context"
	"net/url"
	"testing"
	"time"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_thread_response.xml"

func TestQueryThread(t *testing.T) {
	defer testutils.ActivateMocks()()

	mockURL := constants.ThreadEndpoint + "?id=123"
	testutils.SetupMockResponder(t, mockURL, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	thread, err := Query(context.Background(), client, 123)
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, thread, "Thread should not be nil")

	expected := &ThreadDetail{
		ID:          123,
		NumArticles: 1,
		Link:        "https://boardgamegeek.com/thread/123",
		Subject:     "Example Thread Subject",
		Articles: []Article{
			{
				ID:       456,
				Username: "example_user",
				Link:     "https://boardgamegeek.com/thread/123/article/456#456",
				PostDate: "2023-01-15T10:00:00-05:00",
				EditDate: "2023-01-15T10:00:00-05:00",
				NumEdits: 0,
				Subject:  "Example Article Subject",
				Body:     "This is example content for testing purposes.",
			},
		},
	}

	if diff := cmp.Diff(expected, thread); diff != "" {
		t.Errorf("Thread mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.ThreadEndpoint + "?id=123"

	queryWrapper := func(url string) (*ThreadDetail, error) {
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, 123)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestThreadOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		expected map[string]string
	}{
		{"WithMinArticleID", WithMinArticleID(456), map[string]string{"minarticleid": "456"}},
		{"WithMinArticleDate", WithMinArticleDate(time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)), map[string]string{"minarticledate": "2025-01-15 10:30:00"}},
		{"WithCount", WithCount(50), map[string]string{"count": "50"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			require.NoError(t, tt.option(params))

			for key, expectedValue := range tt.expected {
				assert.Equal(t, expectedValue, params.Get(key))
			}
			assert.Equal(t, len(tt.expected), len(params))
		})
	}
}

func TestThreadOptions_Invalid(t *testing.T) {
	params := url.Values{}
	assert.ErrorIs(t, WithMinArticleID(0)(params), gogeek.ErrInvalidOption)
	assert.ErrorIs(t, WithCount(0)(params), gogeek.ErrInvalidOption)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	result, err := Query(context.Background(), client, 123, WithCount(-1))
	assert.ErrorIs(t, err, gogeek.ErrInvalidOption)
	assert.Nil(t, result)
}
