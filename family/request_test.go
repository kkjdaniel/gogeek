package family

import (
	"context"
	"testing"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_family_response.xml"

func TestQueryFamily(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.FamilyEndpoint + "?id=12&type=" + BoardGameFamily
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	family, err := Query(context.Background(), client, 12, BoardGameFamily)
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

func TestQueryFamily_MultipleIDs(t *testing.T) {
	defer testutils.ActivateMocks()()

	// The mock only answers the URL carrying both IDs, so this pins the
	// comma-delimited format.
	url := constants.FamilyEndpoint + "?id=12,34&type=" + BoardGameFamily
	testutils.SetupMockResponderWithBody(t, url,
		`<items><item type="boardgamefamily" id="12"></item><item type="boardgamefamily" id="34"></item></items>`, 200)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	family, err := Query(context.Background(), client, 12, BoardGameFamily, 34)
	require.NoError(t, err, "Query should accept additional IDs")
	require.Len(t, family.Items, 2, "both families should be returned")
	require.Equal(t, 12, family.Items[0].ID)
	require.Equal(t, 34, family.Items[1].ID)
}

func TestQueryFamily_Error(t *testing.T) {
	testURL := constants.FamilyEndpoint + "?id=12"

	queryWrapper := func(url string) (*Family, error) {
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, 12, BoardGameFamily)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
