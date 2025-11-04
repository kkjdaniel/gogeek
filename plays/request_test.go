package plays

import (

	"github.com/kkjdaniel/gogeek/v2"
	"testing"

	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_plays_response.xml"

func TestQueryPlays(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.PlaysEndpoint + "?username=example_user"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient()
	plays, err := Query(client, "example_user")
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, plays, "Plays should not be nil")

	expected := &Plays{
		Username: "example_user",
		UserID:   123,
		Total:    15,
		Page:     1,
		Plays: []Play{
			{
				ID:         97385232,
				Date:       "2025-04-03",
				Quantity:   1,
				Length:     140,
				Incomplete: 0,
				NoWinStats: 0,
				Location:   "Home",
				Item: PlayItem{
					Name:       "Example Card Game",
					ObjectType: "thing",
					ObjectID:   205637,
					Subtypes: []Subtype{
						{Value: "boardgame"},
					},
				},
				Comments: "Love it! Better than sliced bread.",
				Players: []Player{
					{
						Username: "example_user",
						UserID:   123,
						Name:     "Example User",
						Color:    "Blue",
						Score:    4,
						New:      0,
						Rating:   0,
						Win:      1,
					},
					{
						Username: "",
						UserID:   0,
						Name:     "Player Two",
						Color:    "Red",
						Score:    4,
						New:      0,
						Rating:   0,
						Win:      1,
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, plays); diff != "" {
		t.Errorf("Plays mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.PlaysEndpoint + "?username=example_user"

	queryWrapper := func(url string) (*Plays, error) {
		client := gogeek.NewClient()
		return Query(client, "example_user")
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
