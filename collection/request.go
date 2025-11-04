package collection

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/request"
)

// CollectionOption represents an option for filtering collection queries
type CollectionOption func(params url.Values)

// Query retrieves a user's board game collection from the BoardGameGeek API.
//
// Parameters:
//   - client: A GoGeek client configured with optional authentication
//   - username: A string containing the BGG username whose collection to retrieve
//   - opts: Optional parameters for filtering and customizing the query
//
// Returns:
//   - *Collection: A pointer to a Collection struct containing the user's board game collection
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient()
//	collection, err := collection.Query(client, "exampleuser",
//	    collection.WithOwned(true),
//	    collection.WithStats(),
//	    collection.WithMinRating(7))
func Query(client *gogeek.Client, username string, opts ...CollectionOption) (*Collection, error) {
	params := url.Values{}
	params.Set("username", username)

	// Apply all options
	for _, opt := range opts {
		opt(params)
	}

	queryURL := constants.CollectionEndpoint + "?" + params.Encode()

	var collection Collection
	if err := request.FetchAndUnmarshal(client, queryURL, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// WithVersion adds version info for each item in the collection
func WithVersion() CollectionOption {
	return func(params url.Values) {
		params.Set("version", "1")
	}
}

// WithSubtype specifies which collection type to retrieve
// Valid values: boardgame, boardgameexpansion, boardgameaccessory, rpgitem, rpgissue, videogame
func WithSubtype(subtype string) CollectionOption {
	return func(params url.Values) {
		params.Set("subtype", subtype)
	}
}

// WithExcludeSubtype specifies which subtype to exclude from the results
func WithExcludeSubtype(subtype string) CollectionOption {
	return func(params url.Values) {
		params.Set("excludesubtype", subtype)
	}
}

// WithItemIDs filters collection to specific item IDs
func WithItemIDs(ids ...int) CollectionOption {
	return func(params url.Values) {
		idStrings := make([]string, len(ids))
		for i, id := range ids {
			idStrings[i] = strconv.Itoa(id)
		}
		params.Set("id", strings.Join(idStrings, ","))
	}
}

// WithBrief returns abbreviated results
func WithBrief() CollectionOption {
	return func(params url.Values) {
		params.Set("brief", "1")
	}
}

// WithStats returns expanded rating/ranking info
func WithStats() CollectionOption {
	return func(params url.Values) {
		params.Set("stats", "1")
	}
}

// WithOwned filters for owned games
func WithOwned(owned bool) CollectionOption {
	return func(params url.Values) {
		if owned {
			params.Set("own", "1")
		} else {
			params.Set("own", "0")
		}
	}
}

// WithRated filters for whether an item has been rated
func WithRated(rated bool) CollectionOption {
	return func(params url.Values) {
		if rated {
			params.Set("rated", "1")
		} else {
			params.Set("rated", "0")
		}
	}
}

// WithPlayed filters for whether an item has been played
func WithPlayed(played bool) CollectionOption {
	return func(params url.Values) {
		if played {
			params.Set("played", "1")
		} else {
			params.Set("played", "0")
		}
	}
}

// WithComment filters for items that have been commented
func WithComment(hasComment bool) CollectionOption {
	return func(params url.Values) {
		if hasComment {
			params.Set("comment", "1")
		} else {
			params.Set("comment", "0")
		}
	}
}

// WithTrade filters for items marked for trade
func WithTrade(forTrade bool) CollectionOption {
	return func(params url.Values) {
		if forTrade {
			params.Set("trade", "1")
		} else {
			params.Set("trade", "0")
		}
	}
}

// WithWant filters for items wanted in trade
func WithWant(wanted bool) CollectionOption {
	return func(params url.Values) {
		if wanted {
			params.Set("want", "1")
		} else {
			params.Set("want", "0")
		}
	}
}

// WithWishlist filters for items on the wishlist
func WithWishlist(onWishlist bool) CollectionOption {
	return func(params url.Values) {
		if onWishlist {
			params.Set("wishlist", "1")
		} else {
			params.Set("wishlist", "0")
		}
	}
}

// WithWishlistPriority filters for wishlist priority
// Valid values: 1-5
func WithWishlistPriority(priority int) CollectionOption {
	return func(params url.Values) {
		if priority >= 1 && priority <= 5 {
			params.Set("wishlistpriority", strconv.Itoa(priority))
		}
	}
}

// WithPreordered filters for pre-ordered games
func WithPreordered(preordered bool) CollectionOption {
	return func(params url.Values) {
		if preordered {
			params.Set("preordered", "1")
		} else {
			params.Set("preordered", "0")
		}
	}
}

// WithWantToPlay filters for items marked as wanting to play
func WithWantToPlay(wantToPlay bool) CollectionOption {
	return func(params url.Values) {
		if wantToPlay {
			params.Set("wanttoplay", "1")
		} else {
			params.Set("wanttoplay", "0")
		}
	}
}

// WithWantToBuy filters for items marked as wanting to buy
func WithWantToBuy(wantToBuy bool) CollectionOption {
	return func(params url.Values) {
		if wantToBuy {
			params.Set("wanttobuy", "1")
		} else {
			params.Set("wanttobuy", "0")
		}
	}
}

// WithPrevOwned filters for games marked previously owned
func WithPrevOwned(prevOwned bool) CollectionOption {
	return func(params url.Values) {
		if prevOwned {
			params.Set("prevowned", "1")
		} else {
			params.Set("prevowned", "0")
		}
	}
}

// WithHasParts filters on whether there is a comment in the Has Parts field
func WithHasParts(hasParts bool) CollectionOption {
	return func(params url.Values) {
		if hasParts {
			params.Set("hasparts", "1")
		} else {
			params.Set("hasparts", "0")
		}
	}
}

// WithWantParts filters on whether there is a comment in the Wants Parts field
func WithWantParts(wantParts bool) CollectionOption {
	return func(params url.Values) {
		if wantParts {
			params.Set("wantparts", "1")
		} else {
			params.Set("wantparts", "0")
		}
	}
}

// WithMinRating filters on minimum personal rating assigned
func WithMinRating(rating float64) CollectionOption {
	return func(params url.Values) {
		if rating >= 1 && rating <= 10 {
			params.Set("minrating", fmt.Sprintf("%.1f", rating))
		}
	}
}

// WithMaxRating filters on maximum personal rating assigned
func WithMaxRating(rating float64) CollectionOption {
	return func(params url.Values) {
		if rating >= 1 && rating <= 10 {
			params.Set("rating", fmt.Sprintf("%.1f", rating))
		}
	}
}

// WithMinBGGRating filters on minimum BGG rating
func WithMinBGGRating(rating float64) CollectionOption {
	return func(params url.Values) {
		params.Set("minbggrating", fmt.Sprintf("%.1f", rating))
	}
}

// WithMaxBGGRating filters on maximum BGG rating
func WithMaxBGGRating(rating float64) CollectionOption {
	return func(params url.Values) {
		params.Set("bggrating", fmt.Sprintf("%.1f", rating))
	}
}

// WithMinPlays filters by minimum number of recorded plays
func WithMinPlays(plays int) CollectionOption {
	return func(params url.Values) {
		params.Set("minplays", strconv.Itoa(plays))
	}
}

// WithMaxPlays filters by maximum number of recorded plays
func WithMaxPlays(plays int) CollectionOption {
	return func(params url.Values) {
		params.Set("maxplays", strconv.Itoa(plays))
	}
}

// WithShowPrivate filters to show private collection info
func WithShowPrivate() CollectionOption {
	return func(params url.Values) {
		params.Set("showprivate", "1")
	}
}

// WithCollectionID restricts results to a specific collection ID
func WithCollectionID(collID int) CollectionOption {
	return func(params url.Values) {
		params.Set("collid", strconv.Itoa(collID))
	}
}

// WithModifiedSince restricts results to items modified since date
func WithModifiedSince(date time.Time) CollectionOption {
	return func(params url.Values) {
		params.Set("modifiedsince", date.Format("2006-01-02 15:04:05"))
	}
}
