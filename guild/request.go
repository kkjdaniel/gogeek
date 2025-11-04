package guild

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/request"
)

// Query retrieves detailed information about a specific guild from the BoardGameGeek API.
//
// The function accepts a guild ID and returns a structured representation
// of the guild details including the guild name, category, website, manager,
// description, and location information.
//
// Parameters:
//   - client: A GoGeek client configured with optional authentication
//   - guildID: An integer ID corresponding to a guild in the BGG database
//
// Returns:
//   - *Guild: A pointer to a Guild struct containing the guild information
//   - error: An error if the API request fails or if the response cannot be parsed
//
// Example:
//
//	client := gogeek.NewClient()
//	guild, err := guild.Query(client, 1234)
//	if err != nil {
//	    log.Fatalf("Failed to get guild info: %v", err)
//	}
//	fmt.Printf("Guild name: %s (managed by %s)\n", guild.Name, guild.Manager)
func Query(client *gogeek.Client, guildID int) (*Guild, error) {
	url := fmt.Sprintf(constants.GuildEndpoint+"?id=%d", guildID)

	var guild Guild

	if err := request.FetchAndUnmarshal(client, url, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}
