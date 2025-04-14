package forum

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query retrieves detailed information about a specific forum from the BoardGameGeek API.
//
// The function accepts a forum ID and returns a structured representation
// of the forum details including the forum title, threads within the forum,
// and metadata such as post counts and dates.
//
// Parameters:
//   - id: An integer ID corresponding to a forum in the BGG database
//
// Returns:
//   - *Forum: A pointer to a Forum struct containing the forum information and threads
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	forum, err := forum.Query(1234)
//	if err != nil {
//	    log.Fatalf("Failed to get forum: %v", err)
//	}
//	fmt.Printf("Forum title: %s (contains %d threads)\n", forum.Title, forum.NumThreads)
func Query(id int) (*Forum, error) {
	url := fmt.Sprintf(constants.ForumEndpoint+"?id=%d", id)

	var forumDetail Forum

	if err := request.FetchAndUnmarshal(url, &forumDetail); err != nil {
		return nil, err
	}

	return &forumDetail, nil
}
