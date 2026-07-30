package thing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// ErrNoIDs is returned when no IDs are provided for a query.
var ErrNoIDs = errors.New("no IDs provided")

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
//
// Returns:
//   - *Items: A pointer to an Items struct containing the detailed information for the requested games
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	details, err := thing.Query(context.Background(), client, []int{174430, 167791})
//	if err != nil {
//	    log.Fatalf("Failed to get game details: %v", err)
//	}
//	fmt.Printf("Retrieved details for %d games\n", len(details.Items))
func Query(ctx context.Context, client *gogeek.Client, ids []int) (*Items, error) {
	if len(ids) == 0 {
		return nil, ErrNoIDs
	}

	var all Items
	for start := 0; start < len(ids); start += maxIDsPerRequest {
		batch := ids[start:min(start+maxIDsPerRequest, len(ids))]

		idStrings := make([]string, len(batch))
		for i, id := range batch {
			idStrings[i] = fmt.Sprintf("%d", id)
		}

		url := fmt.Sprintf("%s?id=%s&stats=1", constants.ThingEndpoint, strings.Join(idStrings, ","))

		var page Items
		if err := request.FetchAndUnmarshal(ctx, client, url, &page); err != nil {
			return nil, err
		}
		all.Items = append(all.Items, page.Items...)
	}

	return &all, nil
}
