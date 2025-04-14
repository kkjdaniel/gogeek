package plays

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query retrieves play information for a specific BoardGameGeek user.
//
// The function accepts a BGG username and returns a structured representation
// of the user's play history, including games played, dates, locations,
// player information, and play statistics.
//
// Parameters:
//   - username: A string containing the BGG username whose play history to retrieve
//
// Returns:
//   - *Plays: A pointer to a Plays struct containing the user's play information
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	plays, err := plays.Query("exampleuser")
//	if err != nil {
//	    log.Fatalf("Failed to retrieve plays: %v", err)
//	}
//	fmt.Printf("Found %d plays for user %s\n", plays.Total, plays.Username)
func Query(username string) (*Plays, error) {
	url := fmt.Sprintf(constants.PlaysEndpoint+"?username=%s", username)

	var plays Plays

	if err := request.FetchAndUnmarshal(url, &plays); err != nil {
		return nil, err
	}

	return &plays, nil
}
