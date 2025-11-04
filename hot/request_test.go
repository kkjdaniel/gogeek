package hot

import (

	"github.com/kkjdaniel/gogeek/v2"
	"testing"

	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_hot_response.xml"

func TestQueryHotItems(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.HotEndpoint + "?type=boardgame"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient()
	hotItems, err := Query(client, ItemTypeBoardGame)
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, hotItems, "Hot items should not be nil")

	expected := &HotItems{
		Items: []HotItem{
			{
				ID:   101001,
				Rank: 1,
				Name: ValueString{
					Value: "Example Strategy Game",
				},
				Thumbnail: ValueString{
					Value: "https://example.com/images/game1_thumbnail.jpg",
				},
				YearPublished: ValueInt{
					Value: 2025,
				},
			},
			{
				ID:   101002,
				Rank: 2,
				Name: ValueString{
					Value: "Sample Card Game",
				},
				Thumbnail: ValueString{
					Value: "https://example.com/images/game2_thumbnail.jpg",
				},
				YearPublished: ValueInt{
					Value: 2025,
				},
			},
			{
				ID:   101003,
				Rank: 3,
				Name: ValueString{
					Value: "Generic Board Game",
				},
				Thumbnail: ValueString{
					Value: "https://example.com/images/game3_thumbnail.jpg",
				},
				YearPublished: ValueInt{
					Value: 2025,
				},
			},
		},
	}

	if diff := cmp.Diff(expected, hotItems); diff != "" {
		t.Errorf("Hot items mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.HotEndpoint + "?type=boardgame"

	queryWrapper := func(url string) (*HotItems, error) {
		client := gogeek.NewClient()
		return Query(client, ItemTypeBoardGame)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
