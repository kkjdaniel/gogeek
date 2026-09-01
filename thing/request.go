package thing

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// validThingTypes are the thing types accepted by the type parameter.
var validThingTypes = map[string]bool{
	"boardgame":          true,
	"boardgameaccessory": true,
	"boardgameexpansion": true,
	"rpgitem":            true,
	"rpgissue":           true,
	"videogame":          true,
}

// ErrNoIDs is returned when no IDs are provided for a query.
var ErrNoIDs = errors.New("no IDs provided")

// Option represents an option for customizing thing queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// maxIDsPerRequest is the number of IDs the BGG API accepts in a single
// thing request.
const maxIDsPerRequest = 20

// Query retrieves detailed information about one or more board games from the BoardGameGeek API.
//
// The function accepts a slice of BGG item IDs and returns a structured representation
// of the corresponding board games' details including names, descriptions, categories,
// mechanics, designers, artists, publishers, and various statistics.
//
// The BGG API accepts at most 20 IDs per request, so larger slices are
// transparently split into sequential batches of 20 and the results merged.
// Note that each batch is a separate rate-limited request, so a large ID
// slice takes correspondingly longer.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - ids: A slice of integer IDs corresponding to board game entries in the BGG database
//   - opts: Optional parameters for requesting additional data
//
// Returns:
//   - *Items: A pointer to an Items struct containing the detailed information for the requested games
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	details, err := thing.Query(context.Background(), client, []int{174430, 167791})
//	if err != nil {
//	    log.Fatalf("Failed to get game details: %v", err)
//	}
//	fmt.Printf("Retrieved details for %d games\n", len(details.Items))
func Query(ctx context.Context, client *gogeek.Client, ids []int, opts ...Option) (*Items, error) {
	if len(ids) == 0 {
		return nil, ErrNoIDs
	}

	params := url.Values{}
	params.Set("stats", "1")

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	var all Items
	for start := 0; start < len(ids); start += maxIDsPerRequest {
		batch := ids[start:min(start+maxIDsPerRequest, len(ids))]

		idStrings := make([]string, len(batch))
		for i, id := range batch {
			idStrings[i] = fmt.Sprintf("%d", id)
		}

		// The id list is kept outside the encoded parameters so its commas
		// stay unescaped.
		queryURL := fmt.Sprintf("%s?id=%s&%s", constants.ThingEndpoint, strings.Join(idStrings, ","), params.Encode())

		var page Items
		if err := request.FetchAndUnmarshal(ctx, client, queryURL, &page); err != nil {
			return nil, err
		}
		all.Items = append(all.Items, page.Items...)
	}

	return &all, nil
}

// WithVideos includes the community-submitted videos for each item.
//
// The endpoint returns only the most recent videos, capped at a handful per
// item, while Videos.Total reports how many exist in total.
func WithVideos() Option {
	return func(params url.Values) error {
		params.Set("videos", "1")
		return nil
	}
}

// WithType filters the results to the given thing types.
// Valid types: boardgame, boardgameaccessory, boardgameexpansion, rpgitem, rpgissue, videogame
func WithType(thingTypes ...string) Option {
	return func(params url.Values) error {
		if len(thingTypes) == 0 {
			return fmt.Errorf("%w: at least one thing type is required", gogeek.ErrInvalidOption)
		}
		for _, t := range thingTypes {
			if !validThingTypes[t] {
				return fmt.Errorf("%w: invalid thing type %q", gogeek.ErrInvalidOption, t)
			}
		}
		params.Set("type", strings.Join(thingTypes, ","))
		return nil
	}
}

// WithVersions includes the published versions of each item, such as
// individual printings and localisations.
func WithVersions() Option {
	return func(params url.Values) error {
		params.Set("versions", "1")
		return nil
	}
}

// WithMarketplace includes the current BGG marketplace listings for each item.
func WithMarketplace() Option {
	return func(params url.Values) error {
		params.Set("marketplace", "1")
		return nil
	}
}

// WithComments includes one page of user comments for each item, with the
// commenter's rating when they have given one. It cannot be combined with
// WithRatingComments. Use WithPage and WithPageSize to page through them.
func WithComments() Option {
	return func(params url.Values) error {
		if params.Get("ratingcomments") == "1" {
			return fmt.Errorf("%w: comments and rating comments cannot be requested together", gogeek.ErrInvalidOption)
		}
		params.Set("comments", "1")
		return nil
	}
}

// WithRatingComments includes one page of user ratings for each item, sorted
// by ascending rating and carrying the comment when one was left. It cannot
// be combined with WithComments. Use WithPage and WithPageSize to page
// through them.
func WithRatingComments() Option {
	return func(params url.Values) error {
		if params.Get("comments") == "1" {
			return fmt.Errorf("%w: comments and rating comments cannot be requested together", gogeek.ErrInvalidOption)
		}
		params.Set("ratingcomments", "1")
		return nil
	}
}

// WithPage selects the page of comments or ratings to return. The default is
// page 1.
func WithPage(page int) Option {
	return func(params url.Values) error {
		if page < 1 {
			return fmt.Errorf("%w: page must be at least 1, got %d", gogeek.ErrInvalidOption, page)
		}
		params.Set("page", strconv.Itoa(page))
		return nil
	}
}

// WithPageSize sets the number of comments or ratings per page, between 10
// and 100.
func WithPageSize(size int) Option {
	return func(params url.Values) error {
		if size < 10 || size > 100 {
			return fmt.Errorf("%w: page size must be between 10 and 100, got %d", gogeek.ErrInvalidOption, size)
		}
		params.Set("pagesize", strconv.Itoa(size))
		return nil
	}
}
