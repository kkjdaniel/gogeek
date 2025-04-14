package guild

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/request"
)

// Query retrieves detailed information about a specific guild from the BoardGameGeek API.
//
// The function accepts a guild ID and returns a structured representation
// of the guild details including the guild name, category, website, manager,
// description, and location information.
//
// Parameters:
//   - guildID: An integer ID corresponding to a guild in the BGG database
//
// Returns:
//   - *Guild: A pointer to a Guild struct containing the guild information
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	guild, err := guild.Query(1234)
//	if err != nil {
//	    log.Fatalf("Failed to get guild info: %v", err)
//	}
//	fmt.Printf("Guild name: %s (managed by %s)\n", guild.Name, guild.Manager)
func Query(guildID int) (*Guild, error) {
	url := fmt.Sprintf(constants.GuildEndpoint+"?id=%d", guildID)

	var guild Guild

	if err := request.FetchAndUnmarshal(url, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}
