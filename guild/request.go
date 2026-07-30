package guild

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for customizing guild queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// Query retrieves detailed information about a specific guild from the BoardGameGeek API.
//
// The function accepts a guild ID and returns a structured representation
// of the guild details including the guild name, category, website, manager,
// description, and location information. Use WithMembers to include the
// member roster (25 members per page, paginated with WithPage).
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - guildID: An integer ID corresponding to a guild in the BGG database
//   - opts: Optional parameters for customizing the query
//
// Returns:
//   - *Guild: A pointer to a Guild struct containing the guild information
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	guild, err := guild.Query(context.Background(), client, 1234,
//	    guild.WithMembers())
//	if err != nil {
//	    log.Fatalf("Failed to get guild info: %v", err)
//	}
//	fmt.Printf("Guild name: %s (managed by %s)\n", guild.Name, guild.Manager)
func Query(ctx context.Context, client *gogeek.Client, guildID int, opts ...Option) (*Guild, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(guildID))

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	queryURL := constants.GuildEndpoint + "?" + params.Encode()

	var guild Guild

	if err := request.FetchAndUnmarshal(ctx, client, queryURL, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

// WithMembers includes the guild's member roster in the response
// (25 members per page; see WithPage).
func WithMembers() Option {
	return func(params url.Values) error {
		params.Set("members", "1")
		return nil
	}
}

// WithPage requests the given page of the member roster (25 members per
// page). Only meaningful together with WithMembers.
func WithPage(page int) Option {
	return func(params url.Values) error {
		if page < 1 {
			return fmt.Errorf("%w: page must be at least 1, got %d", gogeek.ErrInvalidOption, page)
		}
		params.Set("page", strconv.Itoa(page))
		return nil
	}
}

// WithSort specifies how the member roster is sorted.
// Valid values: username, date
func WithSort(sort string) Option {
	return func(params url.Values) error {
		if sort != "username" && sort != "date" {
			return fmt.Errorf("%w: invalid sort %q (must be username or date)", gogeek.ErrInvalidOption, sort)
		}
		params.Set("sort", sort)
		return nil
	}
}
