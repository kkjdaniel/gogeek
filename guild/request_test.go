package guild

import (

	"github.com/kkjdaniel/gogeek/v2"
	"testing"

	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_guild_response.xml"

func TestQueryGuild(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.GuildEndpoint + "?id=1234"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	client := gogeek.NewClient()
	guild, err := Query(client, 1234)
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, guild, "Guild should not be nil")

	expected := &Guild{
		ID:          1234,
		Name:        "Example Board Gaming Club",
		Created:     "Sun, 23 May 2021 16:33:41 +0000",
		Category:    "group",
		Website:     "https://www.example.com",
		Manager:     "example_user",
		Description: "This is a sample gaming guild used for testing purposes.",
		Location: Location{
			Addr1:           "Example Community Center",
			Addr2:           "123 Main Street",
			City:            "Anytown",
			StateOrProvince: "State",
			PostalCode:      "12345",
			Country:         "Country",
		},
	}

	if diff := cmp.Diff(expected, guild); diff != "" {
		t.Errorf("Guild mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.GuildEndpoint + "?id=1234"

	queryWrapper := func(url string) (*Guild, error) {
		client := gogeek.NewClient()
		return Query(client, 1234)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}
