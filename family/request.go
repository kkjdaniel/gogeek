package family

import (
	"errors"
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

const (
	RPG             = "rpg"
	RPGPeriodical   = "rpgperiodical"
	BoardGameFamily = "boardgamefamily"
)

var ErrInvalidFamilyType = errors.New("invalid family type")

// Query retrieves detailed information about a specific board game family from the BoardGameGeek API.
//
// The function accepts a family ID and returns a structured representation
// of the family details including the family name, description, and links to games
// within that family.
//
// Parameters:
//   - id: An integer ID corresponding to a board game family in the BGG database
//   - familyType: A string indicating the type of family to query.
//     Must be one of the defined constants: family.RPG, family.RPGPeriodical, or family.BoardGameFamily
//
// Returns:
//   - *Items: A pointer to an Items struct containing the family information
//   - error: An error if the API request fails, if the response cannot be parsed,
//     or if an invalid family type is provided
//
// Example:
//
//	family, err := family.Query(12, family.BoardGameFamily)
//	if err != nil {
//	    log.Fatalf("Failed to get family: %v", err)
//	}
//	fmt.Printf("Family: %s (contains %d games)\n", family.Items[0].Name.Value, len(family.Items[0].Links))
func Query(id int, familyType string) (*Family, error) {
	if !isValidFamilyType(familyType) {
		return nil, fmt.Errorf("%w: %s (must be one of: %s, %s, %s)",
			ErrInvalidFamilyType, familyType, RPG, RPGPeriodical, BoardGameFamily)
	}

	url := fmt.Sprintf("%s?id=%d&type=%s", constants.FamilyEndpoint, id, familyType)

	var familyDetail Family

	if err := request.FetchAndUnmarshal(url, &familyDetail); err != nil {
		return nil, err
	}

	return &familyDetail, nil
}

func isValidFamilyType(familyType string) bool {
	return familyType == RPG || familyType == RPGPeriodical || familyType == BoardGameFamily
}
