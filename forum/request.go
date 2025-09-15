package forum

import (
	"net/url"
	"strconv"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// ForumOption represents an option for customizing forum queries
type ForumOption func(params url.Values)

// Query retrieves detailed information about a specific forum from the BoardGameGeek API.
//
// The function accepts a forum ID and optional parameters, returning a structured representation
// of the forum details including the forum title, threads within the forum,
// and metadata such as post counts and dates.
//
// Parameters:
//   - id: An integer ID corresponding to a forum in the BGG database
//   - opts: Optional parameters for customizing the query (e.g., WithPage for pagination)
//
// Returns:
//   - *Forum: A pointer to a Forum struct containing the forum information and threads
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	// Get first page (default)
//	forum, err := forum.Query(1234)
//	if err != nil {
//	    log.Fatalf("Failed to get forum: %v", err)
//	}
//	fmt.Printf("Forum title: %s (contains %d threads)\n", forum.Title, forum.NumThreads)
//
//	// Get specific page
//	forum, err := forum.Query(1234, forum.WithPage(2))
//	if err != nil {
//	    log.Fatalf("Failed to get forum page 2: %v", err)
//	}
func Query(id int, opts ...ForumOption) (*Forum, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))

	// Apply all options
	for _, opt := range opts {
		opt(params)
	}

	queryURL := constants.ForumEndpoint + "?" + params.Encode()

	var forumDetail Forum

	if err := request.FetchAndUnmarshal(queryURL, &forumDetail); err != nil {
		return nil, err
	}

	return &forumDetail, nil
}

// WithPage specifies which page of threads to retrieve (page size is 50)
// Threads are sorted in order of most recent post.
func WithPage(page int) ForumOption {
	return func(params url.Values) {
		if page > 0 {
			params.Set("page", strconv.Itoa(page))
		}
	}
}
