package thread

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// Option represents an option for customizing thread queries.
// Options validate their arguments and return an error for invalid values.
type Option func(params url.Values) error

// Query retrieves detailed information about a specific thread from the BoardGameGeek API.
//
// The function accepts a thread ID and returns a structured representation
// of the thread details including the thread subject, list of articles posted to the thread,
// and metadata such as post dates and authors.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - threadID: An integer ID corresponding to a thread in the BGG forums
//   - opts: Optional parameters for restricting which articles are returned
//
// Returns:
//   - *ThreadDetail: A pointer to a ThreadDetail struct containing the thread information and articles
//   - error: An error if an option is invalid, the API request fails, or the
//     response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	thread, err := thread.Query(context.Background(), client, 123456)
//	if err != nil {
//	    log.Fatalf("Failed to get thread: %v", err)
//	}
//	fmt.Printf("Thread subject: %s (contains %d articles)\n", thread.Subject, len(thread.Articles))
func Query(ctx context.Context, client *gogeek.Client, threadID int, opts ...Option) (*ThreadDetail, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(threadID))

	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}

	queryURL := constants.ThreadEndpoint + "?" + params.Encode()

	var threadDetail ThreadDetail

	if err := request.FetchAndUnmarshal(ctx, client, queryURL, &threadDetail); err != nil {
		return nil, err
	}

	return &threadDetail, nil
}

// WithMinArticleID restricts results to articles with an ID equal to or
// greater than the given article ID.
func WithMinArticleID(id int) Option {
	return func(params url.Values) error {
		if id < 1 {
			return fmt.Errorf("%w: article id must be at least 1, got %d", gogeek.ErrInvalidOption, id)
		}
		params.Set("minarticleid", strconv.Itoa(id))
		return nil
	}
}

// WithMinArticleDate restricts results to articles posted on or after the
// given date.
func WithMinArticleDate(date time.Time) Option {
	return func(params url.Values) error {
		params.Set("minarticledate", date.Format("2006-01-02 15:04:05"))
		return nil
	}
}

// WithCount limits the number of articles returned.
func WithCount(count int) Option {
	return func(params url.Values) error {
		if count < 1 {
			return fmt.Errorf("%w: count must be at least 1, got %d", gogeek.ErrInvalidOption, count)
		}
		params.Set("count", strconv.Itoa(count))
		return nil
	}
}
