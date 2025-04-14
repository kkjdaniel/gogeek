package thread

import (
	"testing"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_thread_response.xml"

func TestQueryThread(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ThreadEndpoint + "?id=123"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	thread, err := Query(123)
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
