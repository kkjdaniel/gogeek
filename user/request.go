package user

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query retrieves detailed information about a specific user from the BoardGameGeek API.
//
// The function accepts a BGG username and returns a structured representation
// of the user's profile including their basic information, buddies list,
// guild memberships, and top rated items.
//
// Parameters:
//   - username: A string containing the BGG username to retrieve information for
//
// Returns:
//   - *User: A pointer to a User struct containing the user's profile information
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	userProfile, err := user.Query("exampleuser")
//	if err != nil {
//	    log.Fatalf("Failed to retrieve user profile: %v", err)
//	}
//	fmt.Printf("User: %s (member since %s)\n", userProfile.Name, userProfile.YearRegistered)
func Query(username string) (*User, error) {
	url := fmt.Sprintf(constants.UserEndpoint+"?name=%s&buddies=1&guilds=1&top=1", username)

	var user User

	if err := request.FetchAndUnmarshal(url, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
