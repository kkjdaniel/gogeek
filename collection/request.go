package collection

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for filtering collection queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// validSubtypes are the values accepted by the collection endpoint's
// subtype and excludesubtype parameters.
var validSubtypes = map[string]bool{
	"boardgame":          true,
	"boardgameexpansion": true,
	"boardgameaccessory": true,
	"rpgitem":            true,
	"rpgissue":           true,
	"videogame":          true,
}

// Query retrieves a user's board game collection from the BoardGameGeek API.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - username: A string containing the BGG username whose collection to retrieve
//   - opts: Optional parameters for filtering and customizing the query
//
// Returns:
//   - *Collection: A pointer to a Collection struct containing the user's board game collection
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	collection, err := collection.Query(context.Background(), client, "exampleuser",
//	    collection.WithOwned(true),
//	    collection.WithStats(),
//	    collection.WithMinRating(7))
func Query(ctx context.Context, client *gogeek.Client, username string, opts ...Option) (*Collection, error) {
	params := url.Values{}
	params.Set("username", username)

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	queryURL := constants.CollectionEndpoint + "?" + params.Encode()

	var collection Collection
	if err := request.FetchAndUnmarshal(ctx, client, queryURL, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// setFlag encodes a boolean filter as the "1"/"0" form the BGG API expects.
func setFlag(params url.Values, key string, value bool) {
	if value {
		params.Set(key, "1")
	} else {
		params.Set(key, "0")
	}
}

// WithVersion adds version info for each item in the collection
func WithVersion() Option {
	return func(params url.Values) error {
		params.Set("version", "1")
		return nil
	}
}

// WithSubtype specifies which collection type to retrieve
// Valid values: boardgame, boardgameexpansion, boardgameaccessory, rpgitem, rpgissue, videogame
func WithSubtype(subtype string) Option {
	return func(params url.Values) error {
		if !validSubtypes[subtype] {
			return fmt.Errorf("%w: invalid subtype %q", gogeek.ErrInvalidOption, subtype)
		}
		params.Set("subtype", subtype)
		return nil
	}
}

// WithExcludeSubtype specifies which subtype to exclude from the results
// Valid values: boardgame, boardgameexpansion, boardgameaccessory, rpgitem, rpgissue, videogame
func WithExcludeSubtype(subtype string) Option {
	return func(params url.Values) error {
		if !validSubtypes[subtype] {
			return fmt.Errorf("%w: invalid subtype %q", gogeek.ErrInvalidOption, subtype)
		}
		params.Set("excludesubtype", subtype)
		return nil
	}
}

// WithItemIDs filters collection to specific item IDs
func WithItemIDs(ids ...int) Option {
	return func(params url.Values) error {
		if len(ids) == 0 {
			return fmt.Errorf("%w: at least one item ID is required", gogeek.ErrInvalidOption)
		}
		idStrings := make([]string, len(ids))
		for i, id := range ids {
			idStrings[i] = strconv.Itoa(id)
		}
		params.Set("id", strings.Join(idStrings, ","))
		return nil
	}
}

// WithBrief returns abbreviated results
func WithBrief() Option {
	return func(params url.Values) error {
		params.Set("brief", "1")
		return nil
	}
}

// WithStats returns expanded rating/ranking info
func WithStats() Option {
	return func(params url.Values) error {
		params.Set("stats", "1")
		return nil
	}
}

// WithOwned filters for owned games
func WithOwned(owned bool) Option {
	return func(params url.Values) error {
		setFlag(params, "own", owned)
		return nil
	}
}

// WithRated filters for whether an item has been rated
func WithRated(rated bool) Option {
	return func(params url.Values) error {
		setFlag(params, "rated", rated)
		return nil
	}
}

// WithPlayed filters for whether an item has been played
func WithPlayed(played bool) Option {
	return func(params url.Values) error {
		setFlag(params, "played", played)
		return nil
	}
}

// WithComment filters for items that have been commented
func WithComment(hasComment bool) Option {
	return func(params url.Values) error {
		setFlag(params, "comment", hasComment)
		return nil
	}
}

// WithTrade filters for items marked for trade
func WithTrade(forTrade bool) Option {
	return func(params url.Values) error {
		setFlag(params, "trade", forTrade)
		return nil
	}
}

// WithWant filters for items wanted in trade
func WithWant(wanted bool) Option {
	return func(params url.Values) error {
		setFlag(params, "want", wanted)
		return nil
	}
}

// WithWishlist filters for items on the wishlist
func WithWishlist(onWishlist bool) Option {
	return func(params url.Values) error {
		setFlag(params, "wishlist", onWishlist)
		return nil
	}
}

// WithWishlistPriority filters for wishlist priority
// Valid values: 1-5
func WithWishlistPriority(priority int) Option {
	return func(params url.Values) error {
		if priority < 1 || priority > 5 {
			return fmt.Errorf("%w: wishlist priority must be between 1 and 5, got %d", gogeek.ErrInvalidOption, priority)
		}
		params.Set("wishlistpriority", strconv.Itoa(priority))
		return nil
	}
}

// WithPreordered filters for pre-ordered games
func WithPreordered(preordered bool) Option {
	return func(params url.Values) error {
		setFlag(params, "preordered", preordered)
		return nil
	}
}

// WithWantToPlay filters for items marked as wanting to play
func WithWantToPlay(wantToPlay bool) Option {
	return func(params url.Values) error {
		setFlag(params, "wanttoplay", wantToPlay)
		return nil
	}
}

// WithWantToBuy filters for items marked as wanting to buy
func WithWantToBuy(wantToBuy bool) Option {
	return func(params url.Values) error {
		setFlag(params, "wanttobuy", wantToBuy)
		return nil
	}
}

// WithPrevOwned filters for games marked previously owned
func WithPrevOwned(prevOwned bool) Option {
	return func(params url.Values) error {
		setFlag(params, "prevowned", prevOwned)
		return nil
	}
}

// WithHasParts filters on whether there is a comment in the Has Parts field
func WithHasParts(hasParts bool) Option {
	return func(params url.Values) error {
		setFlag(params, "hasparts", hasParts)
		return nil
	}
}

// WithWantParts filters on whether there is a comment in the Wants Parts field
func WithWantParts(wantParts bool) Option {
	return func(params url.Values) error {
		setFlag(params, "wantparts", wantParts)
		return nil
	}
}

// WithMinRating filters on minimum personal rating assigned
// Valid values: 1-10
func WithMinRating(rating float64) Option {
	return func(params url.Values) error {
		if rating < 1 || rating > 10 {
			return fmt.Errorf("%w: rating must be between 1 and 10, got %g", gogeek.ErrInvalidOption, rating)
		}
		params.Set("minrating", fmt.Sprintf("%.1f", rating))
		return nil
	}
}

// WithMaxRating filters on maximum personal rating assigned
// Valid values: 1-10
func WithMaxRating(rating float64) Option {
	return func(params url.Values) error {
		if rating < 1 || rating > 10 {
			return fmt.Errorf("%w: rating must be between 1 and 10, got %g", gogeek.ErrInvalidOption, rating)
		}
		params.Set("rating", fmt.Sprintf("%.1f", rating))
		return nil
	}
}

// WithMinBGGRating filters on minimum BGG rating
// Valid values: 1-10
func WithMinBGGRating(rating float64) Option {
	return func(params url.Values) error {
		if rating < 1 || rating > 10 {
			return fmt.Errorf("%w: rating must be between 1 and 10, got %g", gogeek.ErrInvalidOption, rating)
		}
		params.Set("minbggrating", fmt.Sprintf("%.1f", rating))
		return nil
	}
}

// WithMaxBGGRating filters on maximum BGG rating
// Valid values: 1-10
func WithMaxBGGRating(rating float64) Option {
	return func(params url.Values) error {
		if rating < 1 || rating > 10 {
			return fmt.Errorf("%w: rating must be between 1 and 10, got %g", gogeek.ErrInvalidOption, rating)
		}
		params.Set("bggrating", fmt.Sprintf("%.1f", rating))
		return nil
	}
}

// WithMinPlays filters by minimum number of recorded plays
func WithMinPlays(plays int) Option {
	return func(params url.Values) error {
		if plays < 0 {
			return fmt.Errorf("%w: plays must be non-negative, got %d", gogeek.ErrInvalidOption, plays)
		}
		params.Set("minplays", strconv.Itoa(plays))
		return nil
	}
}

// WithMaxPlays filters by maximum number of recorded plays
func WithMaxPlays(plays int) Option {
	return func(params url.Values) error {
		if plays < 0 {
			return fmt.Errorf("%w: plays must be non-negative, got %d", gogeek.ErrInvalidOption, plays)
		}
		params.Set("maxplays", strconv.Itoa(plays))
		return nil
	}
}

// WithShowPrivate filters to show private collection info
func WithShowPrivate() Option {
	return func(params url.Values) error {
		params.Set("showprivate", "1")
		return nil
	}
}

// WithCollectionID restricts results to a specific collection ID
func WithCollectionID(collID int) Option {
	return func(params url.Values) error {
		params.Set("collid", strconv.Itoa(collID))
		return nil
	}
}

// WithModifiedSince restricts results to items modified since date
func WithModifiedSince(date time.Time) Option {
	return func(params url.Values) error {
		params.Set("modifiedsince", date.Format("2006-01-02 15:04:05"))
		return nil
	}
}
