package thread

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query retrieves detailed information about a specific thread from the BoardGameGeek API.
//
// The function accepts a thread ID and returns a structured representation
// of the thread details including the thread subject, list of articles posted to the thread,
// and metadata such as post dates and authors.
//
// Parameters:
//   - threadID: An integer ID corresponding to a thread in the BGG forums
//
// Returns:
//   - *ThreadDetail: A pointer to a ThreadDetail struct containing the thread information and articles
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	thread, err := thread.Query(123456)
//	if err != nil {
//	    log.Fatalf("Failed to get thread: %v", err)
//	}
//	fmt.Printf("Thread subject: %s (contains %d articles)\n", thread.Subject, len(thread.Articles))
func Query(threadID int) (*ThreadDetail, error) {
	url := fmt.Sprintf(constants.ThreadEndpoint+"?id=%d", threadID)

	var threadDetail ThreadDetail

	if err := request.FetchAndUnmarshal(url, &threadDetail); err != nil {
		return nil, err
	}

	return &threadDetail, nil
}
