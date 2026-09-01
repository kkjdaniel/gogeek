package family

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

const (
	// RPG is the family type for RPG families.
	RPG = "rpg"
	// RPGPeriodical is the family type for RPG periodical families.
	RPGPeriodical = "rpgperiodical"
	// BoardGameFamily is the family type for board game families.
	BoardGameFamily = "boardgamefamily"
)

// ErrInvalidFamilyType is returned when the provided family type is not one of
// RPG, RPGPeriodical, or BoardGameFamily.
var ErrInvalidFamilyType = errors.New("invalid family type")

// Query retrieves detailed information about one or more board game families from the BoardGameGeek API.
//
// The function accepts family IDs and returns a structured representation
// of the family details including the family name, description, and links to games
// within each family.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - id: An integer ID corresponding to a board game family in the BGG database
//   - familyType: A string indicating the type of family to query.
//     Must be one of the defined constants: family.RPG, family.RPGPeriodical, or family.BoardGameFamily
//   - moreIDs: Optional additional family IDs to retrieve in the same request
//
// Returns:
//   - *Items: A pointer to an Items struct containing the family information
//   - error: An error if the API request fails, if the response cannot be parsed,
//     or if an invalid family type is provided
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	family, err := family.Query(context.Background(), client, 12, family.BoardGameFamily)
//	if err != nil {
//	    log.Fatalf("Failed to get family: %v", err)
//	}
//	fmt.Printf("Family: %s (contains %d games)\n", family.Items[0].Name.Value, len(family.Items[0].Links))
func Query(ctx context.Context, client *gogeek.Client, id int, familyType string, moreIDs ...int) (*Family, error) {
	if !isValidFamilyType(familyType) {
		return nil, fmt.Errorf("%w: %s (must be one of: %s, %s, %s)",
			ErrInvalidFamilyType, familyType, RPG, RPGPeriodical, BoardGameFamily)
	}

	idStrings := make([]string, 0, 1+len(moreIDs))
	for _, v := range append([]int{id}, moreIDs...) {
		idStrings = append(idStrings, strconv.Itoa(v))
	}

	url := fmt.Sprintf("%s?id=%s&type=%s", constants.FamilyEndpoint, strings.Join(idStrings, ","), familyType)

	var familyDetail Family

	if err := request.FetchAndUnmarshal(ctx, client, url, &familyDetail); err != nil {
		return nil, err
	}

	return &familyDetail, nil
}

func isValidFamilyType(familyType string) bool {
	return familyType == RPG || familyType == RPGPeriodical || familyType == BoardGameFamily
}
