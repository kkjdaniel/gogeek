package user

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for customizing user queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// validDomains are the values accepted by the user endpoint's domain parameter.
var validDomains = map[string]bool{
	"boardgame": true,
	"rpg":       true,
	"videogame": true,
}

// Query retrieves detailed information about a specific user from the BoardGameGeek API.
//
// The function accepts a BGG username and returns a structured representation
// of the user's profile. By default the buddies list, guild memberships, and
// top rated items are included; use WithoutBuddies, WithoutGuilds, and
// WithoutTop to omit them. Usernames with spaces or special characters are
// automatically URL-encoded.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - username: A string containing the BGG username to retrieve information for
//   - opts: Optional parameters for customizing the query
//
// Returns:
//   - *User: A pointer to a User struct containing the user's profile information
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	userProfile, err := user.Query(context.Background(), client, "example user")
//	if err != nil {
//	    log.Fatalf("Failed to retrieve user profile: %v", err)
//	}
//	fmt.Printf("User: %s (member since %s)\n", userProfile.Name, userProfile.YearRegistered)
func Query(ctx context.Context, client *gogeek.Client, username string, opts ...Option) (*User, error) {
	params := url.Values{}
	params.Set("name", username)
	params.Set("buddies", "1")
	params.Set("guilds", "1")
	params.Set("top", "1")

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	requestURL := constants.UserEndpoint + "?" + params.Encode()

	var user User

	if err := request.FetchAndUnmarshal(ctx, client, requestURL, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// WithoutBuddies omits the user's buddies list from the response.
func WithoutBuddies() Option {
	return func(params url.Values) error {
		params.Del("buddies")
		return nil
	}
}

// WithoutGuilds omits the user's guild memberships from the response.
func WithoutGuilds() Option {
	return func(params url.Values) error {
		params.Del("guilds")
		return nil
	}
}

// WithoutTop omits the user's top rated items from the response.
func WithoutTop() Option {
	return func(params url.Values) error {
		params.Del("top")
		return nil
	}
}

// WithHot includes the user's hot list in the response.
func WithHot() Option {
	return func(params url.Values) error {
		params.Set("hot", "1")
		return nil
	}
}

// WithDomain specifies the domain for the user's hot and top lists.
// Valid values: boardgame (default), rpg, videogame
func WithDomain(domain string) Option {
	return func(params url.Values) error {
		if !validDomains[domain] {
			return fmt.Errorf("%w: invalid domain %q (must be boardgame, rpg, or videogame)", gogeek.ErrInvalidOption, domain)
		}
		params.Set("domain", domain)
		return nil
	}
}

// WithPage requests the given page of the buddies and guilds lists
// (100 entries per page).
func WithPage(page int) Option {
	return func(params url.Values) error {
		if page < 1 {
			return fmt.Errorf("%w: page must be at least 1, got %d", gogeek.ErrInvalidOption, page)
		}
		params.Set("page", strconv.Itoa(page))
		return nil
	}
}
