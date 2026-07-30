package hot

import (
	"context"
	"fmt"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/request"
)

// ItemType identifies which category of trending items to retrieve from
// the hotness endpoint.
type ItemType string

// The item types accepted by the hot endpoint.
const (
	ItemTypeBoardGame        ItemType = "boardgame"
	ItemTypeRPG              ItemType = "rpg"
	ItemTypeVideoGame        ItemType = "videogame"
	ItemTypeBoardGamePerson  ItemType = "boardgameperson"
	ItemTypeRPGPerson        ItemType = "rpgperson"
	ItemTypeBoardGameCompany ItemType = "boardgamecompany"
	ItemTypeRPGCompany       ItemType = "rpgcompany"
	ItemTypeVideoGameCompany ItemType = "videogamecompany"
)

// validItemTypes is the set of item types accepted by the hot endpoint.
var validItemTypes = map[ItemType]bool{
	ItemTypeBoardGame:        true,
	ItemTypeRPG:              true,
	ItemTypeVideoGame:        true,
	ItemTypeBoardGamePerson:  true,
	ItemTypeRPGPerson:        true,
	ItemTypeBoardGameCompany: true,
	ItemTypeRPGCompany:       true,
	ItemTypeVideoGameCompany: true,
}

// Query retrieves the current "hotness" list from the BoardGameGeek API for a specific item type.
//
// The function accepts an item type parameter and returns a structured representation
// of the items currently trending on BoardGameGeek, including their ranks, names,
// and basic metadata.
//
// Parameters:
//   - ctx: A context that can cancel or time-bound the request
//   - client: A GoGeek client configured with authentication
//   - itemType: An ItemType constant specifying which category of items to retrieve
//     (e.g. ItemTypeBoardGame, ItemTypeRPG)
//
// Returns:
//   - *HotItems: A pointer to a HotItems struct containing the list of trending items
//   - error: An error if the item type is invalid, the API request fails, or
//     the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
//	hotGames, err := hot.Query(context.Background(), client, hot.ItemTypeBoardGame)
//	if err != nil {
//	    log.Fatalf("Failed to retrieve hot games: %v", err)
//	}
//	fmt.Printf("Retrieved %d hot games. #1 is %s\n", len(hotGames.Items), hotGames.Items[0].Name.Value)
func Query(ctx context.Context, client *gogeek.Client, itemType ItemType) (*HotItems, error) {
	if !validItemTypes[itemType] {
		return nil, fmt.Errorf("%w: invalid item type %q", gogeek.ErrInvalidOption, itemType)
	}

	url := fmt.Sprintf(constants.HotEndpoint+"?type=%s", itemType)

	var hotItems HotItems
	if err := request.FetchAndUnmarshal(ctx, client, url, &hotItems); err != nil {
		return nil, err
	}

	return &hotItems, nil
}
