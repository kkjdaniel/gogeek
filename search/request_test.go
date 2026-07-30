package search

import (
	"context"
	"net/url"
	"testing"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_search_response.xml"

func TestQuerySearch(t *testing.T) {
	defer testutils.ActivateMocks()()

	mockURL := constants.SearchEndpoint + "?query=test"
	testutils.SetupMockResponder(t, mockURL, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	results, err := Query(context.Background(), client, "test")
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, results, "Search results should not be nil")

	expected := &SearchResults{
		Total: 4,
		Items: []SearchResult{
			{
				ID:   134277,
				Type: "boardgame",
				Name: Name{
					Type:  "alternate",
					Value: "Example Board Game Expansion",
				},
				YearPublished: YearPublishedTag{
					Value: 2012,
				},
			},
			{
				ID:   110308,
				Type: "boardgame",
				Name: Name{
					Type:  "primary",
					Value: "Sample Strategy Game",
				},
				YearPublished: YearPublishedTag{
					Value: 2011,
				},
			},
			{
				ID:   123386,
				Type: "boardgame",
				Name: Name{
					Type:  "primary",
					Value: "Generic Board Game",
				},
				YearPublished: YearPublishedTag{
					Value: 2012,
				},
			},
			{
				ID:   5824,
				Type: "boardgame",
				Name: Name{
					Type:  "alternate",
					Value: "Test Family Game",
				},
				YearPublished: YearPublishedTag{
					Value: 2003,
				},
			},
		},
	}

	if diff := cmp.Diff(expected, results); diff != "" {
		t.Errorf("Search results mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.SearchEndpoint + "?query=test"

	queryWrapper := func(url string) (*SearchResults, error) {
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, "test")
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestSearchOptions(t *testing.T) {
	t.Run("WithExact", func(t *testing.T) {
		params := url.Values{}
		require.NoError(t, WithExact()(params))
		assert.Equal(t, "1", params.Get("exact"))
	})

	t.Run("WithType single", func(t *testing.T) {
		params := url.Values{}
		require.NoError(t, WithType("boardgame")(params))
		assert.Equal(t, "boardgame", params.Get("type"))
	})

	t.Run("WithType multiple", func(t *testing.T) {
		params := url.Values{}
		require.NoError(t, WithType("boardgame", "boardgameexpansion")(params))
		assert.Equal(t, "boardgame,boardgameexpansion", params.Get("type"))
	})

	t.Run("WithType invalid", func(t *testing.T) {
		params := url.Values{}
		assert.ErrorIs(t, WithType("bogus")(params), gogeek.ErrInvalidOption)
		assert.ErrorIs(t, WithType()(params), gogeek.ErrInvalidOption)
	})

	t.Run("Query surfaces option error", func(t *testing.T) {
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		result, err := Query(context.Background(), client, "test", WithType("bogus"))
		assert.ErrorIs(t, err, gogeek.ErrInvalidOption)
		assert.Nil(t, result)
	})
}
