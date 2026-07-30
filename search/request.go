package search

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for customizing search queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// validItemTypes are the item types accepted by the search endpoint's
// type parameter.
var validItemTypes = map[string]bool{
	"boardgame":          true,
	"boardgameaccessory": true,
	"boardgameexpansion": true,
	"rpgitem":            true,
	"videogame":          true,
}

// Query searches for items in the BoardGameGeek database using a text query.
//
// The function accepts a search query string and returns a structured representation
// of the search results including item IDs, names, and publication years.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - query: A string containing the search terms to find matching items
//   - opts: Optional parameters for customizing the query (e.g., WithExact, WithType)
//
// Returns:
//   - *SearchResults: A pointer to a SearchResults struct containing the search results
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//
//	// Regular search
//	results, err := search.Query(context.Background(), client, "catan")
//
//	// Exact match search restricted to board games
//	exact, err := search.Query(context.Background(), client, "catan",
//	    search.WithExact(),
//	    search.WithType("boardgame"))
func Query(ctx context.Context, client *gogeek.Client, query string, opts ...Option) (*SearchResults, error) {
	params := url.Values{}
	params.Set("query", query)

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	requestURL := constants.SearchEndpoint + "?" + params.Encode()

	var searchResults SearchResults
	if err := request.FetchAndUnmarshal(ctx, client, requestURL, &searchResults); err != nil {
		return nil, err
	}

	return &searchResults, nil
}

// WithExact limits results to exact matches of the query.
func WithExact() Option {
	return func(params url.Values) error {
		params.Set("exact", "1")
		return nil
	}
}

// WithType restricts results to the given item types.
// Valid types: boardgame, boardgameaccessory, boardgameexpansion, rpgitem, videogame
func WithType(itemTypes ...string) Option {
	return func(params url.Values) error {
		if len(itemTypes) == 0 {
			return fmt.Errorf("%w: at least one item type is required", gogeek.ErrInvalidOption)
		}
		for _, t := range itemTypes {
			if !validItemTypes[t] {
				return fmt.Errorf("%w: invalid item type %q", gogeek.ErrInvalidOption, t)
			}
		}
		params.Set("type", strings.Join(itemTypes, ","))
		return nil
	}
}
