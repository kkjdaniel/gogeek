package collection

import (
	"testing"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_collection_response.xml"

func TestQueryCollection(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.CollectionEndpoint + "?username=testuser"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	collection, err := Query("testuser")
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, collection, "Collection should not be nil")

	expected := &Collection{
		TotalItems: 3,
		PubDate:    "Fri, 04 Apr 2025 11:41:47 +0000",
		Items: []CollectionItem{
			{
				ObjectType:    "thing",
				ObjectID:      101001,
				Subtype:       "boardgame",
				CollectionID:  201001,
				Name:          "Example Strategy Card Game",
				YearPublished: 2022,
				Image:         "https://example.com/images/game1_full.jpg",
				Thumbnail:     "https://example.com/images/game1_thumb.jpg",
				Status: ItemStatus{
					Own:          0,
					PrevOwned:    0,
					ForTrade:     0,
					Want:         0,
					WantToPlay:   0,
					WantToBuy:    0,
					Wishlist:     0,
					Preordered:   0,
					LastModified: "2025-03-12 08:02:55",
				},
				NumPlays: 0,
				Comment:  "This is an example comment about a board game.",
			},
			{
				ObjectType:    "thing",
				ObjectID:      101002,
				Subtype:       "boardgame",
				CollectionID:  201002,
				Name:          "Sample Economic Game",
				YearPublished: 2020,
				Image:         "https://example.com/images/game2_full.jpg",
				Thumbnail:     "https://example.com/images/game2_thumb.jpg",
				Status: ItemStatus{
					Own:          1,
					PrevOwned:    0,
					ForTrade:     0,
					Want:         0,
					WantToPlay:   0,
					WantToBuy:    0,
					Wishlist:     0,
					Preordered:   0,
					LastModified: "2025-03-16 10:32:55",
				},
				NumPlays: 3,
			},
			{
				ObjectType:    "thing",
				ObjectID:      101003,
				Subtype:       "boardgame",
				CollectionID:  201003,
				Name:          "Generic Family Game",
				YearPublished: 2018,
				Image:         "https://example.com/images/game3_full.jpg",
				Thumbnail:     "https://example.com/images/game3_thumb.jpg",
				Status: ItemStatus{
					Own:          0,
					PrevOwned:    0,
					ForTrade:     0,
					Want:         0,
					WantToPlay:   1,
					WantToBuy:    0,
					Wishlist:     1,
					Preordered:   0,
					LastModified: "2025-02-12 02:58:00",
				},
				NumPlays: 0,
			},
		},
	}

	if diff := cmp.Diff(expected, collection); diff != "" {
		t.Errorf("Collection mismatch (-want +got):\n%s", diff)
	}
}
