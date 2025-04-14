package search

import (
	"testing"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_search_response.xml"

func TestQuerySearch(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.SearchEndpoint + "?query=test"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	results, err := Query("test")
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
