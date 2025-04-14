package hot

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

type ItemType string

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

// Query retrieves the current "hotness" list from the BoardGameGeek API for a specific item type.
//
// The function accepts an item type parameter and returns a structured representation
// of the items currently trending on BoardGameGeek, including their ranks, names,
// and basic metadata.
//
// Parameters:
//   - itemType: An ItemType constant specifying which category of items to retrieve
//     (e.g. ItemTypeBoardGame, ItemTypeRPG)
//
// Returns:
//   - *HotItems: A pointer to a HotItems struct containing the list of trending items
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	hotGames, err := hot.Query(hot.ItemTypeBoardGame)
//	if err != nil {
//	    log.Fatalf("Failed to retrieve hot games: %v", err)
//	}
//	fmt.Printf("Retrieved %d hot games. #1 is %s\n", len(hotGames.Items), hotGames.Items[0].Name.Value)
func Query(itemType ItemType) (*HotItems, error) {
	url := fmt.Sprintf(constants.HotEndpoint+"?type=%s", itemType)

	var hotItems HotItems
	if err := request.FetchAndUnmarshal(url, &hotItems); err != nil {
		return nil, err
	}

	return &hotItems, nil
}
