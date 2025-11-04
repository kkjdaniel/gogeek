package family

import (
	"testing"

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_family_response.xml"

func TestQueryFamily(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.FamilyEndpoint + "?id=12&type=" + BoardGameFamily
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient()
	family, err := Query(client, 12, BoardGameFamily)
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, family, "Family should not be nil")

	expected := &Family{
		Items: []Item{
			{
				Type:      "boardgamefamily",
				ID:        12,
				Thumbnail: "https://example.com/images/family_thumbnail.jpg",
				Image:     "https://example.com/images/family_full.jpg",
				Name: Name{
					Type:      "primary",
					SortIndex: 1,
					Value:     "Sample Game Series",
				},
				Description: "This is an example description for a board game family.",
				Links: []Link{
					{Type: "boardgamefamily", ID: 101, Value: "Sample Game 1", Inbound: true},
					{Type: "boardgamefamily", ID: 102, Value: "Sample Game 2", Inbound: true},
					{Type: "boardgamefamily", ID: 103, Value: "Sample Game 3", Inbound: true},
					{Type: "boardgamefamily", ID: 104, Value: "Sample Game 4", Inbound: true},
					{Type: "boardgamefamily", ID: 105, Value: "Sample Game 5", Inbound: true},
					{Type: "boardgamefamily", ID: 106, Value: "Sample Game 6", Inbound: true},
					{Type: "boardgamefamily", ID: 107, Value: "Sample Game 7", Inbound: true},
					{Type: "boardgamefamily", ID: 108, Value: "Sample Game 8", Inbound: true},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, family); diff != "" {
		t.Errorf("Family mismatch (-want +got):\n%s", diff)
	}
}

func TestQueryFamily_Error(t *testing.T) {
	testURL := constants.FamilyEndpoint + "?id=12"

	queryWrapper := func(url string) (*Family, error) {
		client := gogeek.NewClient()
		return Query(client, 12, BoardGameFamily)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
