package forum

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for customizing forum queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// Query retrieves detailed information about a specific forum from the BoardGameGeek API.
//
// The function accepts a forum ID and optional parameters, returning a structured representation
// of the forum details including the forum title, threads within the forum,
// and metadata such as post counts and dates.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - id: An integer ID corresponding to a forum in the BGG database
//   - opts: Optional parameters for customizing the query (e.g., WithPage for pagination)
//
// Returns:
//   - *Forum: A pointer to a Forum struct containing the forum information and threads
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	// Get first page (default)
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	forum, err := forum.Query(context.Background(), client, 1234)
//	if err != nil {
//	    log.Fatalf("Failed to get forum: %v", err)
//	}
//	fmt.Printf("Forum title: %s (contains %d threads)\n", forum.Title, forum.NumThreads)
//
//	// Get specific page
//	forum, err = forum.Query(context.Background(), client, 1234, forum.WithPage(2))
//	if err != nil {
//	    log.Fatalf("Failed to get forum page 2: %v", err)
//	}
func Query(ctx context.Context, client *gogeek.Client, id int, opts ...Option) (*Forum, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	queryURL := constants.ForumEndpoint + "?" + params.Encode()

	var forumDetail Forum

	if err := request.FetchAndUnmarshal(ctx, client, queryURL, &forumDetail); err != nil {
		return nil, err
	}

	return &forumDetail, nil
}

// WithPage specifies which page of threads to retrieve (page size is 50).
// Threads are sorted in order of most recent post.
func WithPage(page int) Option {
	return func(params url.Values) error {
		if page < 1 {
			return fmt.Errorf("%w: page must be at least 1, got %d", gogeek.ErrInvalidOption, page)
		}
		params.Set("page", strconv.Itoa(page))
		return nil
	}
}
