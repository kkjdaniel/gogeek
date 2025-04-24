package thing

import (
	"testing"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_thing_response.xml"

func TestQueryThing(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ThingEndpoint + "?id=9"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	thing, err := Query([]int{9})
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, thing, "Thing should not be nil")

	expected := &Items{
		Items: []Item{
			{
				Type: "boardgame",
				ID:   9,
				Name: []Name{
					{Type: "primary", SortIndex: 1, Value: "Example Game"},
					{Type: "alternate", SortIndex: 2, Value: "Sample Game"},
					{Type: "alternate", SortIndex: 3, Value: "Test Game"},
					{Type: "alternate", SortIndex: 4, Value: "Demo Game"},
				},
				Description:   "This is an example description for a board game.",
				Thumbnail:     "https://example.com/images/game_thumbnail.jpg",
				Image:         "https://example.com/images/game_full.jpg",
				YearPublished: YearPublished{Value: 2000},
				MinPlayers:    Value{Value: 2},
				MaxPlayers:    Value{Value: 4},
				PlayingTime:   Value{Value: 90},
				MinPlayTime:   Value{Value: 90},
				MaxPlayTime:   Value{Value: 90},
				MinAge:        Value{Value: 10},
				Links: []Link{
					{Type: "boardgamecategory", ID: 1001, Value: "Strategy"},
					{Type: "boardgamemechanic", ID: 2001, Value: "Area Control"},
					{Type: "boardgamemechanic", ID: 2002, Value: "Tile Placement"},
					{Type: "boardgamefamily", ID: 3001, Value: "Game Series: Example Games"},
					{Type: "boardgamedesigner", ID: 4001, Value: "Designer One"},
					{Type: "boardgamedesigner", ID: 4002, Value: "Designer Two"},
					{Type: "boardgameartist", ID: 5001, Value: "Artist Name"},
					{Type: "boardgamepublisher", ID: 6001, Value: "Publisher One"},
					{Type: "boardgamepublisher", ID: 6002, Value: "Publisher Two"},
					{Type: "boardgamepublisher", ID: 6003, Value: "Publisher Three"},
				},
				Polls: []Poll{
					{
						Name:       "suggested_numplayers",
						Title:      "User Suggested Number of Players",
						TotalVotes: 60,
						Results: []PollResult{
							{
								NumPlayers: "1",
								Values:     []ResultValue{{Value: "Best", NumVotes: 0}, {Value: "Recommended", NumVotes: 0}, {Value: "Not Recommended", NumVotes: 30}},
							},
							{
								NumPlayers: "2",
								Values:     []ResultValue{{Value: "Best", NumVotes: 10}, {Value: "Recommended", NumVotes: 20}, {Value: "Not Recommended", NumVotes: 5}},
							},
							{
								NumPlayers: "3",
								Values:     []ResultValue{{Value: "Best", NumVotes: 25}, {Value: "Recommended", NumVotes: 20}, {Value: "Not Recommended", NumVotes: 0}},
							},
							{
								NumPlayers: "4",
								Values:     []ResultValue{{Value: "Best", NumVotes: 15}, {Value: "Recommended", NumVotes: 25}, {Value: "Not Recommended", NumVotes: 5}},
							},
							{
								NumPlayers: "4+",
								Values:     []ResultValue{{Value: "Best", NumVotes: 0}, {Value: "Recommended", NumVotes: 0}, {Value: "Not Recommended", NumVotes: 30}},
							},
						},
					},
					{
						Name:       "suggested_playerage",
						Title:      "User Suggested Player Age",
						TotalVotes: 10,
						Results: []PollResult{
							{
								Values: []ResultValue{
									{Value: "2", NumVotes: 0},
									{Value: "3", NumVotes: 0},
									{Value: "4", NumVotes: 0},
									{Value: "5", NumVotes: 0},
									{Value: "6", NumVotes: 0},
									{Value: "8", NumVotes: 2},
									{Value: "10", NumVotes: 4},
									{Value: "12", NumVotes: 3},
									{Value: "14", NumVotes: 1},
									{Value: "16", NumVotes: 0},
									{Value: "18", NumVotes: 0},
									{Value: "21 and up", NumVotes: 0},
								},
							},
						},
					},
					{
						Name:       "language_dependence",
						Title:      "Language Dependence",
						TotalVotes: 9,
						Results: []PollResult{
							{
								Values: []ResultValue{
									{Value: "No necessary in-game text", NumVotes: 8},
									{Value: "Some necessary text - easily memorized or small crib sheet", NumVotes: 1},
									{Value: "Moderate in-game text - needs crib sheet or paste ups", NumVotes: 0},
									{Value: "Extensive use of text - massive conversion needed to be playable", NumVotes: 0},
									{Value: "Unplayable in another language", NumVotes: 0},
								},
							},
						},
					},
				},
				PollSummaries: []PollSummary{
					{
						Name:  "suggested_numplayers",
						Title: "User Suggested Number of Players",
						Results: []SummaryItem{
							{Name: "bestwith", Value: "Best with 3 players"},
							{Name: "recommmendedwith", Value: "Recommended with 2–4 players"},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, thing); diff != "" {
		t.Errorf("Thing mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.ThingEndpoint + "?id=9"

	queryWrapper := func(url string) (*Items, error) {
		return Query([]int{9})
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
