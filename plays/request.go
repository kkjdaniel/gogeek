package plays

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/request"
)

// Query retrieves play information for a specific BoardGameGeek user.
//
// The function accepts a BGG username and returns a structured representation
// of the user's play history, including games played, dates, locations,
// player information, and play statistics.
//
// Parameters:
//   - client: A GoGeek client configured with optional authentication
//   - username: A string containing the BGG username whose play history to retrieve
//
// Returns:
//   - *Plays: A pointer to a Plays struct containing the user's play information
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient()
//	plays, err := plays.Query(client, "exampleuser")
//	if err != nil {
//	    log.Fatalf("Failed to retrieve plays: %v", err)
//	}
//	fmt.Printf("Found %d plays for user %s\n", plays.Total, plays.Username)
func Query(client *gogeek.Client, username string) (*Plays, error) {
	url := fmt.Sprintf(constants.PlaysEndpoint+"?username=%s", username)

	var plays Plays

	if err := request.FetchAndUnmarshal(client, url, &plays); err != nil {
		return nil, err
	}

	return &plays, nil
}
