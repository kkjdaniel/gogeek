package plays

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

const mockDataFileValid = "testdata/valid_plays_response.xml"

func TestQueryPlays(t *testing.T) {
	defer testutils.ActivateMocks()()

	mockURL := constants.PlaysEndpoint + "?username=example_user"
	testutils.SetupMockResponder(t, mockURL, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	plays, err := Query(context.Background(), client, "example_user")
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
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, "example_user")
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestPlaysOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		expected map[string]string
	}{
		{"WithPage", WithPage(2), map[string]string{"page": "2"}},
		{"WithMinDate", WithMinDate(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)), map[string]string{"mindate": "2025-01-15"}},
		{"WithMaxDate", WithMaxDate(time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)), map[string]string{"maxdate": "2025-06-30"}},
		{"WithID", WithID(13), map[string]string{"id": "13"}},
		{"WithType", WithType("thing"), map[string]string{"type": "thing"}},
		{"WithSubtype", WithSubtype("boardgame"), map[string]string{"subtype": "boardgame"}},
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

func TestPlaysOptions_Invalid(t *testing.T) {
	params := url.Values{}
	assert.ErrorIs(t, WithPage(0)(params), gogeek.ErrInvalidOption)
	assert.ErrorIs(t, WithID(0)(params), gogeek.ErrInvalidOption)
	assert.ErrorIs(t, WithType("bogus")(params), gogeek.ErrInvalidOption)
	assert.ErrorIs(t, WithSubtype("bogus")(params), gogeek.ErrInvalidOption)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	result, err := Query(context.Background(), client, "example_user", WithPage(0))
	assert.ErrorIs(t, err, gogeek.ErrInvalidOption)
	assert.Nil(t, result)
}
