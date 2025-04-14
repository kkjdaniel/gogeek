package search

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query searches for board games in the BoardGameGeek database using a text query.
//
// The function accepts a search query string and returns a structured representation
// of the search results including game IDs, names, and publication years.
//
// Parameters:
//   - query: A string containing the search terms to find matching board games
//
// Returns:
//   - *SearchResults: A pointer to a SearchResults struct containing the search results
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	results, err := search.Query("catan")
//	if err != nil {
//	    log.Fatalf("Failed to search for games: %v", err)
//	}
//	fmt.Printf("Found %d results for 'catan'\n", results.Total)
func Query(query string) (*SearchResults, error) {
	url := fmt.Sprintf(constants.SearchEndpoint+"?query=%s", query)

	var searchResults SearchResults

	if err := request.FetchAndUnmarshal(url, &searchResults); err != nil {
		return nil, err
	}

	return &searchResults, nil
}
