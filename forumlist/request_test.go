package forumlist

import (
	"testing"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const (
	mockDataFileValid       = "testdata/valid_forumlist_response.xml"
	mockDataFileValidFamily = "testdata/valid_family_forumlist_response.xml"
)

func TestFetchForumList(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ForumListEndpoint + "?id=174430&type=thing"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	forumList, err := Query(174430, Thing)
	require.NoError(t, err, "FetchForumList should not return an error")
	require.NotNil(t, forumList, "ForumList should not be nil")

	expected := &ForumList{
		Type: "thing",
		ID:   174430,
		Forums: []Forum{
			{
				ID:           0,
				GroupID:      0,
				Title:        "Reviews",
				NoPosting:    0,
				Description:  "Post your game reviews in this forum.",
				NumThreads:   84,
				NumPosts:     139,
				LastPostDate: "Wed, 03 Jan 2024 16:27:59 +0000",
			},
			{
				ID:           1,
				GroupID:      0,
				Title:        "Sessions",
				NoPosting:    0,
				Description:  "Post your session reports here.",
				NumThreads:   453,
				NumPosts:     1098,
				LastPostDate: "Mon, 20 Jan 2025 04:19:24 +0000",
			},
			{
				ID:           2,
				GroupID:      0,
				Title:        "General",
				NoPosting:    0,
				Description:  "Discuss this game and organize play sessions.",
				NumThreads:   2459,
				NumPosts:     23986,
				LastPostDate: "Mon, 27 Jan 2025 17:18:42 +0000",
			},
			{
				ID:           65,
				GroupID:      0,
				Title:        "Rules",
				NoPosting:    0,
				Description:  "Post your rules questions here.",
				NumThreads:   1072,
				NumPosts:     6456,
				LastPostDate: "Mon, 27 Jan 2025 17:18:31 +0000",
			},
			{
				ID:           67,
				GroupID:      1,
				Title:        "Play By Forum",
				NoPosting:    0,
				Description:  "PBF Games of Gloomhaven",
				NumThreads:   216,
				NumPosts:     90089,
				LastPostDate: "Mon, 27 Jan 2025 20:56:59 +0000",
			},
			{
				ID:           69,
				GroupID:      0,
				Title:        "Variants",
				NoPosting:    0,
				Description:  "Post your variants here.",
				NumThreads:   146,
				NumPosts:     835,
				LastPostDate: "Sun, 26 Jan 2025 06:18:42 +0000",
			},
		},
	}

	if diff := cmp.Diff(expected, forumList); diff != "" {
		t.Errorf("ForumList mismatch (-want +got):\n%s", diff)
	}
}

func TestFetchForumList_Family(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ForumListEndpoint + "?id=12&type=family"
	testutils.SetupMockResponder(t, url, mockDataFileValidFamily)

	forumList, err := Query(12, Family)
	require.NoError(t, err, "FetchForumList should not return an error for family type")
	require.NotNil(t, forumList, "ForumList should not be nil")

	expected := &ForumList{
		Type: "family",
		ID:   12,
		Forums: []Forum{
			{
				ID:           1001,
				GroupID:      0,
				Title:        "General Discussion",
				NoPosting:    0,
				Description:  "General discussion about this game family.",
				NumThreads:   25,
				NumPosts:     187,
				LastPostDate: "Mon, 15 Jan 2025 09:30:00 +0000",
			},
			{
				ID:           1002,
				GroupID:      0,
				Title:        "Comparisons",
				NoPosting:    0,
				Description:  "Compare games within this family.",
				NumThreads:   12,
				NumPosts:     89,
				LastPostDate: "Thu, 10 Jan 2025 14:22:30 +0000",
			},
			{
				ID:           1003,
				GroupID:      0,
				Title:        "Recommendations",
				NoPosting:    0,
				Description:  "Get recommendations for games in this family.",
				NumThreads:   8,
				NumPosts:     42,
				LastPostDate: "Sun, 07 Jan 2025 18:45:00 +0000",
			},
		},
	}

	if diff := cmp.Diff(expected, forumList); diff != "" {
		t.Errorf("ForumList mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_InvalidType(t *testing.T) {
	_, err := Query(174430, "invalid")
	require.Error(t, err, "Query should return an error for invalid type")
	require.ErrorIs(t, err, ErrInvalidForumListType, "Error should be ErrInvalidForumListType")
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.ForumListEndpoint + "?id=174430&type=thing"

	queryWrapper := func(url string) (*ForumList, error) {
		return Query(174430, Thing)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}