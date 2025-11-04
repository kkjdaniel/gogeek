package forum

import (
	"testing"

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_forum_response.xml"

func TestFetchForum(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ForumEndpoint + "?id=123"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient()
	forum, err := Query(client, 123)
	require.NoError(t, err, "FetchForum should not return an error")
	require.NotNil(t, forum, "Forum should not be nil")

	expected := &Forum{
		ID:           123,
		Title:        "Example Forum",
		NumThreads:   3,
		NumPosts:     499,
		LastPostDate: "Thu, 01 Jan 2025 00:00:00 +0000",
		NoPosting:    0,
		Threads: []Thread{
			{
				ID:           1234,
				Subject:      "Example Thread Topic 1",
				Author:       "example_user1",
				NumArticles:  1,
				PostDate:     "Tue, 28 Jan 2025 04:50:26 +0000",
				LastPostDate: "Tue, 28 Jan 2025 04:50:26 +0000",
			},
			{
				ID:           1235,
				Subject:      "Example Thread Topic 2",
				Author:       "example_user2",
				NumArticles:  1,
				PostDate:     "Thu, 23 Jan 2025 18:37:40 +0000",
				LastPostDate: "Thu, 23 Jan 2025 18:37:40 +0000",
			},
			{
				ID:           1236,
				Subject:      "Example Thread Topic 3",
				Author:       "example_user3",
				NumArticles:  45,
				PostDate:     "Sun, 10 Jan 2025 23:43:41 +0000",
				LastPostDate: "Tue, 14 Jan 2025 00:45:23 +0000",
			},
		},
	}

	if diff := cmp.Diff(expected, forum); diff != "" {
		t.Errorf("Forum mismatch (-want +got):\n%s", diff)
	}
}

func TestFetchForum_WithPage(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ForumEndpoint + "?id=123&page=2"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient()
	forum, err := Query(client, 123, WithPage(2))
	require.NoError(t, err, "FetchForum with page should not return an error")
	require.NotNil(t, forum, "Forum should not be nil")
	require.Equal(t, 123, forum.ID)
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.ForumEndpoint + "?id=123"

	queryWrapper := func(url string) (*Forum, error) {
		client := gogeek.NewClient()
		return Query(client, 123)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
