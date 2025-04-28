package search

import (
	"net/url"

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
//   - exact: Optional boolean parameter to enable exact matching
//
// Returns:
//   - *SearchResults: A pointer to a SearchResults struct containing the search results
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	// Regular search
//	results, err := search.Query("catan")
//
//	// Exact match search
//	exactResults, err := search.Query("catan", true)
func Query(query string, exact ...bool) (*SearchResults, error) {
	params := url.Values{}
	params.Set("query", query)

	if len(exact) > 0 && exact[0] {
		params.Set("exact", "1")
	}

	requestURL := constants.SearchEndpoint + "?" + params.Encode()

	var searchResults SearchResults
	if err := request.FetchAndUnmarshal(requestURL, &searchResults); err != nil {
		return nil, err
	}

	return &searchResults, nil
}
