package thing

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const (
	mockDataFileValid       = "testdata/valid_thing_response.xml"
	mockDataFileValidVideos = "testdata/valid_thing_videos_response.xml"
)

func TestQueryThing(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ThingEndpoint + "?id=9&stats=1"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	thing, err := Query(context.Background(), client, []int{9})
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
				YearPublished: IntValue{Value: 2000},
				MinPlayers:    IntValue{Value: 2},
				MaxPlayers:    IntValue{Value: 4},
				PlayingTime:   IntValue{Value: 90},
				MinPlayTime:   IntValue{Value: 90},
				MaxPlayTime:   IntValue{Value: 90},
				MinAge:        IntValue{Value: 10},
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
				Statistics: &Statistics{
					UsersRated:   IntValue{Value: 3929},
					Average:      FloatValue{Value: 7.28028},
					BayesAverage: FloatValue{Value: 6.59011},
					Ranks: []Rank{
						{
							Type:         "subtype",
							ID:           1,
							Name:         "boardgame",
							Friendly:     "Board Game Rank",
							Value:        "1071",
							BayesAverage: "6.59011",
						},
						{
							Type:         "family",
							ID:           5498,
							Name:         "partygames",
							Friendly:     "Party Game Rank",
							Value:        "49",
							BayesAverage: "6.85912",
						},
						{
							Type:         "family",
							ID:           5499,
							Name:         "familygames",
							Friendly:     "Family Game Rank",
							Value:        "276",
							BayesAverage: "6.72714",
						},
					},
					StdDev:        FloatValue{Value: 1.41125},
					Median:        IntValue{Value: 0},
					Owned:         IntValue{Value: 8727},
					Trading:       IntValue{Value: 41},
					Wanting:       IntValue{Value: 220},
					Wishing:       IntValue{Value: 1653},
					NumComments:   IntValue{Value: 718},
					NumWeights:    IntValue{Value: 91},
					AverageWeight: FloatValue{Value: 1.0989},
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
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, []int{9})
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestQuery_NoIDs(t *testing.T) {
	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	result, err := Query(context.Background(), client, nil)

	require.ErrorIs(t, err, ErrNoIDs)
	require.Nil(t, result)
}

func TestQuery_AutoChunksBeyond20IDs(t *testing.T) {
	defer testutils.ActivateMocks()()

	// 25 IDs should be fetched as one batch of 20 and one batch of 5.
	ids := make([]int, 25)
	for i := range ids {
		ids[i] = i + 1
	}

	makeBody := func(ids []int) string {
		var sb strings.Builder
		sb.WriteString(`<items>`)
		for _, id := range ids {
			fmt.Fprintf(&sb, `<item type="boardgame" id="%d"></item>`, id)
		}
		sb.WriteString(`</items>`)
		return sb.String()
	}

	firstIDs := make([]string, 20)
	for i := range firstIDs {
		firstIDs[i] = fmt.Sprintf("%d", i+1)
	}
	secondIDs := make([]string, 5)
	for i := range secondIDs {
		secondIDs[i] = fmt.Sprintf("%d", i+21)
	}

	firstURL := fmt.Sprintf("%s?id=%s&stats=1", constants.ThingEndpoint, strings.Join(firstIDs, ","))
	secondURL := fmt.Sprintf("%s?id=%s&stats=1", constants.ThingEndpoint, strings.Join(secondIDs, ","))

	testutils.SetupMockResponderWithBody(t, firstURL, makeBody(ids[:20]), 200)
	testutils.SetupMockResponderWithBody(t, secondURL, makeBody(ids[20:]), 200)

	client := gogeek.NewClient(gogeek.APIKey("test-key"), gogeek.WithRateLimit(1000))
	result, err := Query(context.Background(), client, ids)

	require.NoError(t, err, "Query should transparently chunk >20 IDs")
	require.Len(t, result.Items, 25, "All items from both batches should be merged")
	require.Equal(t, 1, result.Items[0].ID)
	require.Equal(t, 25, result.Items[24].ID)
}

func TestQueryThing_WithVideos(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.ThingEndpoint + "?id=9&stats=1&videos=1"
	testutils.SetupMockResponder(t, url, mockDataFileValidVideos)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	thing, err := Query(context.Background(), client, []int{9}, WithVideos())
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, thing, "Thing should not be nil")
	require.Len(t, thing.Items, 1)

	videos := thing.Items[0].Videos
	require.NotNil(t, videos, "Videos should be populated when WithVideos is used")

	expected := &Videos{
		Total: 4,
		Videos: []Video{
			{
				ID:       501,
				Title:    "Example Game Review",
				Category: "review",
				Language: "English",
				Link:     "https://www.youtube.com/watch?v=example1",
				Username: "reviewer_one",
				UserID:   1001,
				PostDate: "2020-01-15T12:00:00-06:00",
			},
			{
				ID:       502,
				Title:    "Comment jouer a Example Game",
				Category: "instructional",
				Language: "French",
				Link:     "https://www.youtube.com/watch?v=example2",
				Username: "joueur_deux",
				UserID:   1002,
				PostDate: "2021-06-02T09:30:00-05:00",
			},
			{
				ID:       503,
				Title:    "Example Game Unboxing",
				Category: "unboxing",
				Language: "English",
				Link:     "https://www.youtube.com/watch?v=example3",
				Username: "unboxer_three",
				UserID:   1003,
				PostDate: "2022-11-20T18:45:00-06:00",
			},
			{
				ID:       504,
				Title:    "Example Game Session Report",
				Category: "session",
				Language: "German",
				Link:     "https://www.youtube.com/watch?v=example4",
				Username: "spieler_vier",
				UserID:   1004,
				PostDate: "2023-03-08T07:15:00-06:00",
			},
		},
	}

	if diff := cmp.Diff(expected, videos); diff != "" {
		t.Errorf("Videos mismatch (-want +got):\n%s", diff)
	}
}

func TestQueryThing_WithoutVideosIsUnchanged(t *testing.T) {
	defer testutils.ActivateMocks()()

	// Query is part of the public API, so a call with no options must still
	// produce exactly the URL it always has, without a videos parameter. The
	// mock only answers that URL, so any drift fails this test.
	url := constants.ThingEndpoint + "?id=9&stats=1"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	thing, err := Query(context.Background(), client, []int{9})
	require.NoError(t, err)
	require.Len(t, thing.Items, 1)
	require.Nil(t, thing.Items[0].Videos, "Videos should stay nil when not requested")
}

func TestQueryThing_OptionErrorIsReturned(t *testing.T) {
	failing := func(url.Values) error { return errors.New("boom") }

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	result, err := Query(context.Background(), client, []int{9}, failing)

	require.Error(t, err, "an option error should abort the query")
	require.Nil(t, result)
}

func TestQueryThing_VideosAcrossBatches(t *testing.T) {
	defer testutils.ActivateMocks()()

	// Options must be applied to every batch, not only the first.
	ids := make([]int, 25)
	for i := range ids {
		ids[i] = i + 1
	}

	makeBody := func(ids []int) string {
		var sb strings.Builder
		sb.WriteString(`<items>`)
		for _, id := range ids {
			fmt.Fprintf(&sb, `<item type="boardgame" id="%d"><videos total="1">`+
				`<video id="%d" title="v" category="review" language="English" link="l" username="u" userid="1" postdate="p"/>`+
				`</videos></item>`, id, id)
		}
		sb.WriteString(`</items>`)
		return sb.String()
	}

	idStrings := func(ids []int) string {
		parts := make([]string, len(ids))
		for i, id := range ids {
			parts[i] = fmt.Sprintf("%d", id)
		}
		return strings.Join(parts, ",")
	}

	firstURL := fmt.Sprintf("%s?id=%s&stats=1&videos=1", constants.ThingEndpoint, idStrings(ids[:20]))
	secondURL := fmt.Sprintf("%s?id=%s&stats=1&videos=1", constants.ThingEndpoint, idStrings(ids[20:]))

	testutils.SetupMockResponderWithBody(t, firstURL, makeBody(ids[:20]), 200)
	testutils.SetupMockResponderWithBody(t, secondURL, makeBody(ids[20:]), 200)

	client := gogeek.NewClient(gogeek.APIKey("test-key"), gogeek.WithRateLimit(1000))
	result, err := Query(context.Background(), client, ids, WithVideos())

	require.NoError(t, err)
	require.Len(t, result.Items, 25)
	require.NotNil(t, result.Items[24].Videos, "the second batch should carry videos too")
	require.Equal(t, 1, result.Items[24].Videos.Total)
}
