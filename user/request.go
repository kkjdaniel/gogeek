package user

import (
	"fmt"
	"net/url" // Add this import for URL encoding

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/request"
)

// Query retrieves detailed information about a specific user from the BoardGameGeek API.
//
// The function accepts a BGG username and returns a structured representation
// of the user's profile including their basic information, buddies list,
// guild memberships, and top rated items. Usernames with spaces or special
// characters are automatically URL-encoded.
//
// Parameters:
//   - client: A GoGeek client configured with optional authentication
//   - username: A string containing the BGG username to retrieve information for
//
// Returns:
//   - *User: A pointer to a User struct containing the user's profile information
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient()
//	userProfile, err := user.Query(client, "example user")
//	if err != nil {
//	    log.Fatalf("Failed to retrieve user profile: %v", err)
//	}
//	fmt.Printf("User: %s (member since %s)\n", userProfile.Name, userProfile.YearRegistered)
func Query(client *gogeek.Client, username string) (*User, error) {
	escapedUsername := url.QueryEscape(username)

	requestURL := fmt.Sprintf("%s?name=%s&buddies=1&guilds=1&top=1",
		constants.UserEndpoint, escapedUsername)

	var user User

	if err := request.FetchAndUnmarshal(client, requestURL, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
