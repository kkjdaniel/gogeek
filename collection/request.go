package collection

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query retrieves a user's board game collection from the BoardGameGeek API.
//
// The function accepts a BGG username and returns a structured representation
// of the user's collection, including owned games, wishlist items, and related
// metadata such as play counts and acquisition dates.
//
// Parameters:
//   - username: A string containing the BGG username whose collection to retrieve
//
// Returns:
//   - *Collection: A pointer to a Collection struct containing the user's board game collection
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	collection, err := collection.Query("exampleuser")
//	if err != nil {
//	    log.Fatalf("Failed to retrieve collection: %v", err)
//	}
//	fmt.Printf("Found %d items in %s's collection\n", len(collection.Items), "exampleuser")
func Query(username string) (*Collection, error) {
	url := fmt.Sprintf(constants.CollectionEndpoint+"?username=%s", username)

	var collection Collection

	if err := request.FetchAndUnmarshal(url, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}
