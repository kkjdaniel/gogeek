package collection

import (
	"net/url"
	"testing"
	"time"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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

func TestQuery_Error(t *testing.T) {
	testURL := constants.CollectionEndpoint + "?username=testuser"

	queryWrapper := func(url string) (*Collection, error) {
		return Query("testuser")
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestCollectionOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   CollectionOption
		expected map[string]string
	}{
		{"WithVersion", WithVersion(), map[string]string{"version": "1"}},
		{"WithSubtype", WithSubtype("boardgame"), map[string]string{"subtype": "boardgame"}},
		{"WithExcludeSubtype", WithExcludeSubtype("rpgitem"), map[string]string{"excludesubtype": "rpgitem"}},
		{"WithItemIDs", WithItemIDs(123, 456), map[string]string{"id": "123,456"}},
		{"WithBrief", WithBrief(), map[string]string{"brief": "1"}},
		{"WithStats", WithStats(), map[string]string{"stats": "1"}},
		{"WithOwnedTrue", WithOwned(true), map[string]string{"own": "1"}},
		{"WithOwnedFalse", WithOwned(false), map[string]string{"own": "0"}},
		{"WithRatedTrue", WithRated(true), map[string]string{"rated": "1"}},
		{"WithRatedFalse", WithRated(false), map[string]string{"rated": "0"}},
		{"WithPlayedTrue", WithPlayed(true), map[string]string{"played": "1"}},
		{"WithPlayedFalse", WithPlayed(false), map[string]string{"played": "0"}},
		{"WithCommentTrue", WithComment(true), map[string]string{"comment": "1"}},
		{"WithCommentFalse", WithComment(false), map[string]string{"comment": "0"}},
		{"WithTradeTrue", WithTrade(true), map[string]string{"trade": "1"}},
		{"WithTradeFalse", WithTrade(false), map[string]string{"trade": "0"}},
		{"WithWantTrue", WithWant(true), map[string]string{"want": "1"}},
		{"WithWantFalse", WithWant(false), map[string]string{"want": "0"}},
		{"WithWishlistTrue", WithWishlist(true), map[string]string{"wishlist": "1"}},
		{"WithWishlistFalse", WithWishlist(false), map[string]string{"wishlist": "0"}},
		{"WithWishlistPriority", WithWishlistPriority(3), map[string]string{"wishlistpriority": "3"}},
		{"WithPreorderedTrue", WithPreordered(true), map[string]string{"preordered": "1"}},
		{"WithPreorderedFalse", WithPreordered(false), map[string]string{"preordered": "0"}},
		{"WithWantToPlayTrue", WithWantToPlay(true), map[string]string{"wanttoplay": "1"}},
		{"WithWantToPlayFalse", WithWantToPlay(false), map[string]string{"wanttoplay": "0"}},
		{"WithWantToBuyTrue", WithWantToBuy(true), map[string]string{"wanttobuy": "1"}},
		{"WithWantToBuyFalse", WithWantToBuy(false), map[string]string{"wanttobuy": "0"}},
		{"WithPrevOwnedTrue", WithPrevOwned(true), map[string]string{"prevowned": "1"}},
		{"WithPrevOwnedFalse", WithPrevOwned(false), map[string]string{"prevowned": "0"}},
		{"WithHasPartsTrue", WithHasParts(true), map[string]string{"hasparts": "1"}},
		{"WithHasPartsFalse", WithHasParts(false), map[string]string{"hasparts": "0"}},
		{"WithWantPartsTrue", WithWantParts(true), map[string]string{"wantparts": "1"}},
		{"WithWantPartsFalse", WithWantParts(false), map[string]string{"wantparts": "0"}},
		{"WithMinRating", WithMinRating(7.5), map[string]string{"minrating": "7.5"}},
		{"WithMaxRating", WithMaxRating(9.0), map[string]string{"rating": "9.0"}},
		{"WithMinBGGRating", WithMinBGGRating(8.0), map[string]string{"minbggrating": "8.0"}},
		{"WithMaxBGGRating", WithMaxBGGRating(9.5), map[string]string{"bggrating": "9.5"}},
		{"WithMinPlays", WithMinPlays(5), map[string]string{"minplays": "5"}},
		{"WithMaxPlays", WithMaxPlays(10), map[string]string{"maxplays": "10"}},
		{"WithShowPrivate", WithShowPrivate(), map[string]string{"showprivate": "1"}},
		{"WithCollectionID", WithCollectionID(12345), map[string]string{"collid": "12345"}},
	}

	t.Run("WithModifiedSince", func(t *testing.T) {
		params := url.Values{}
		testDate := time.Date(2025, 4, 1, 12, 30, 0, 0, time.UTC)
		WithModifiedSince(testDate)(params)
		assert.Equal(t, "2025-04-01 12:30:00", params.Get("modifiedsince"))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			tt.option(params)

			for key, expectedValue := range tt.expected {
				assert.Equal(t, expectedValue, params.Get(key),
					"Parameter %s should be %s", key, expectedValue)
			}

			assert.Equal(t, len(tt.expected), len(params),
				"Should have exactly the expected number of parameters")
		})
	}
}

func TestInvalidParameterValues(t *testing.T) {
	t.Run("WishlistPriority invalid bounds", func(t *testing.T) {
		for _, invalid := range []int{0, 6, -1, 10} {
			params := url.Values{}
			WithWishlistPriority(invalid)(params)
			assert.Empty(t, params.Get("wishlistpriority"),
				"Should not set parameter for invalid value %d", invalid)
		}
	})

	t.Run("Rating invalid bounds", func(t *testing.T) {
		for _, invalid := range []float64{0.5, 10.5, -1.0, 11.0} {
			params := url.Values{}
			WithMinRating(invalid)(params)
			assert.Empty(t, params.Get("minrating"),
				"Should not set parameter for invalid value %f", invalid)

			params = url.Values{}
			WithMaxRating(invalid)(params)
			assert.Empty(t, params.Get("rating"),
				"Should not set parameter for invalid value %f", invalid)
		}
	})
}

func TestMultipleOptions(t *testing.T) {
	params := url.Values{}

	options := []CollectionOption{
		WithOwned(true),
		WithStats(),
		WithMinRating(7.0),
		WithPlayed(true),
	}

	for _, opt := range options {
		opt(params)
	}

	expected := map[string]string{
		"own":       "1",
		"stats":     "1",
		"minrating": "7.0",
		"played":    "1",
	}

	for key, value := range expected {
		assert.Equal(t, value, params.Get(key))
	}

	assert.Equal(t, len(expected), len(params))
}
