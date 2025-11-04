package thing

import (
	"fmt"
	"strings"

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/request"
)

var (
	// ErrTooManyIDs is returned when more than 20 IDs are provided for a query.
	ErrTooManyIDs = fmt.Errorf("too many IDs provided, maximum is 20")
	// ErrNoIDs is returned when no IDs are provided for a query.
	ErrNoIDs = fmt.Errorf("no IDs provided")
)

// Query retrieves detailed information about one or more board games from the BoardGameGeek API.
//
// The function accepts a slice of BGG item IDs and returns a structured representation
// of the corresponding board games' details including names, descriptions, categories,
// mechanics, designers, artists, publishers, and various statistics.
//
// Parameters:
//   - client: A GoGeek client configured with optional authentication
//   - ids: A slice of integer IDs corresponding to board game entries in the BGG database
//
// Returns:
//   - *Items: A pointer to an Items struct containing the detailed information for the requested games
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient()
//	details, err := thing.Query(client, []int{174430, 167791})
//	if err != nil {
//	    log.Fatalf("Failed to get game details: %v", err)
//	}
//	fmt.Printf("Retrieved details for %d games\n", len(details.Items))
func Query(client *gogeek.Client, ids []int) (*Items, error) {
	if len(ids) == 0 {
		return nil, ErrNoIDs
	}

	if len(ids) > 20 {
		return nil, ErrTooManyIDs
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = fmt.Sprintf("%d", id)
	}
	idParam := strings.Join(idStrings, ",")

	url := fmt.Sprintf("%s?id=%s&stats=1", constants.ThingEndpoint, idParam)

	var thing Items
	if err := request.FetchAndUnmarshal(client, url, &thing); err != nil {
		return nil, err
	}

	return &thing, nil
}
