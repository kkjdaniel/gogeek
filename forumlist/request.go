package forumlist

import (
	"context"
	"errors"
	"fmt"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

const (
	// Thing is the forum list type for things (e.g. board games).
	Thing = "thing"
	// Family is the forum list type for families.
	Family = "family"
)

// ErrInvalidForumListType is returned when the provided forum list type is not
// one of Thing or Family.
var ErrInvalidForumListType = errors.New("invalid forum list type")

// Query retrieves a list of forums for a particular type/id from the BoardGameGeek API.
//
// The function accepts an ID and a forum list type, returning a structured representation
// of the forums associated with that particular thing or family, including forum titles,
// descriptions, and metadata such as thread and post counts.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - id: An integer ID corresponding to a thing or family in the BGG database
//   - forumListType: A string indicating the type of entry to query.
//     Must be one of the defined constants: forumlist.Thing or forumlist.Family
//
// Returns:
//   - *ForumList: A pointer to a ForumList struct containing the forums information
//   - error: An error if the API request fails, if the response cannot be parsed,
//     or if an invalid forum list type is provided
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	forums, err := forumlist.Query(context.Background(), client, 174430, forumlist.Thing)
//	if err != nil {
//	    log.Fatalf("Failed to get forum list: %v", err)
//	}
//	fmt.Printf("Found %d forums for this game\n", len(forums.Forums))
func Query(ctx context.Context, client *gogeek.Client, id int, forumListType string) (*ForumList, error) {
	if !isValidForumListType(forumListType) {
		return nil, fmt.Errorf("%w: %s (must be one of: %s, %s)",
			ErrInvalidForumListType, forumListType, Thing, Family)
	}

	url := fmt.Sprintf("%s?id=%d&type=%s", constants.ForumListEndpoint, id, forumListType)

	var forumList ForumList

	if err := request.FetchAndUnmarshal(ctx, client, url, &forumList); err != nil {
		return nil, err
	}

	return &forumList, nil
}

func isValidForumListType(forumListType string) bool {
	return forumListType == Thing || forumListType == Family
}
