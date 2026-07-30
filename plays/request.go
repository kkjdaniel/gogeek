package plays

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for filtering and paginating play queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// validTypes are the values accepted by the plays endpoint's type parameter.
var validTypes = map[string]bool{
	"thing":  true,
	"family": true,
}

// validSubtypes are the values accepted by the plays endpoint's subtype parameter.
var validSubtypes = map[string]bool{
	"boardgame":               true,
	"boardgameexpansion":      true,
	"boardgameaccessory":      true,
	"boardgameintegration":    true,
	"boardgamecompilation":    true,
	"boardgameimplementation": true,
	"rpg":                     true,
	"rpgitem":                 true,
	"videogame":               true,
}

// Query retrieves play information for a specific BoardGameGeek user.
//
// The function accepts a BGG username and returns a structured representation
// of the user's play history, including games played, dates, locations,
// player information, and play statistics.
//
// BGG returns at most 100 plays per page; use WithPage to fetch further
// pages when Total exceeds the number of plays returned.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - username: A string containing the BGG username whose play history to retrieve
//   - opts: Optional parameters for filtering and paginating the query
//
// Returns:
//   - *Plays: A pointer to a Plays struct containing the user's play information
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	plays, err := plays.Query(context.Background(), client, "exampleuser",
//	    plays.WithPage(2))
//	if err != nil {
//	    log.Fatalf("Failed to retrieve plays: %v", err)
//	}
//	fmt.Printf("Found %d plays for user %s\n", plays.Total, plays.Username)
func Query(ctx context.Context, client *gogeek.Client, username string, opts ...Option) (*Plays, error) {
	params := url.Values{}
	params.Set("username", username)

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	queryURL := constants.PlaysEndpoint + "?" + params.Encode()

	var plays Plays

	if err := request.FetchAndUnmarshal(ctx, client, queryURL, &plays); err != nil {
		return nil, err
	}

	return &plays, nil
}

// WithPage requests the given page of plays (100 plays per page).
func WithPage(page int) Option {
	return func(params url.Values) error {
		if page < 1 {
			return fmt.Errorf("%w: page must be at least 1, got %d", gogeek.ErrInvalidOption, page)
		}
		params.Set("page", strconv.Itoa(page))
		return nil
	}
}

// WithMinDate restricts results to plays on or after the given date.
func WithMinDate(date time.Time) Option {
	return func(params url.Values) error {
		params.Set("mindate", date.Format("2006-01-02"))
		return nil
	}
}

// WithMaxDate restricts results to plays on or before the given date.
func WithMaxDate(date time.Time) Option {
	return func(params url.Values) error {
		params.Set("maxdate", date.Format("2006-01-02"))
		return nil
	}
}

// WithID restricts results to plays of the item with the given BGG ID.
func WithID(id int) Option {
	return func(params url.Values) error {
		if id < 1 {
			return fmt.Errorf("%w: id must be at least 1, got %d", gogeek.ErrInvalidOption, id)
		}
		params.Set("id", strconv.Itoa(id))
		return nil
	}
}

// WithType restricts results to plays of the given item type.
// Valid types: thing, family
func WithType(itemType string) Option {
	return func(params url.Values) error {
		if !validTypes[itemType] {
			return fmt.Errorf("%w: invalid type %q (must be thing or family)", gogeek.ErrInvalidOption, itemType)
		}
		params.Set("type", itemType)
		return nil
	}
}

// WithSubtype restricts results to plays of the given item subtype.
// Valid subtypes: boardgame, boardgameexpansion, boardgameaccessory,
// boardgameintegration, boardgamecompilation, boardgameimplementation,
// rpg, rpgitem, videogame
func WithSubtype(subtype string) Option {
	return func(params url.Values) error {
		if !validSubtypes[subtype] {
			return fmt.Errorf("%w: invalid subtype %q", gogeek.ErrInvalidOption, subtype)
		}
		params.Set("subtype", subtype)
		return nil
	}
}
