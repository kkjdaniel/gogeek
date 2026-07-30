package guild

import (
	"context"
	"net/url"
	"testing"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/internal/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_guild_response.xml"

func TestQueryGuild(t *testing.T) {
	defer testutils.ActivateMocks()()

	mockURL := constants.GuildEndpoint + "?id=1234"
	testutils.SetupMockResponder(t, mockURL, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	guild, err := Query(context.Background(), client, 1234)
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
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, 1234)
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestGuildOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		expected map[string]string
	}{
		{"WithMembers", WithMembers(), map[string]string{"members": "1"}},
		{"WithPage", WithPage(3), map[string]string{"page": "3"}},
		{"WithSortUsername", WithSort("username"), map[string]string{"sort": "username"}},
		{"WithSortDate", WithSort("date"), map[string]string{"sort": "date"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			require.NoError(t, tt.option(params))

			for key, expectedValue := range tt.expected {
				assert.Equal(t, expectedValue, params.Get(key))
			}
			assert.Equal(t, len(tt.expected), len(params))
		})
	}
}

func TestGuildOptions_Invalid(t *testing.T) {
	params := url.Values{}
	assert.ErrorIs(t, WithPage(0)(params), gogeek.ErrInvalidOption)
	assert.ErrorIs(t, WithSort("bogus")(params), gogeek.ErrInvalidOption)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	result, err := Query(context.Background(), client, 1234, WithSort("bogus"))
	assert.ErrorIs(t, err, gogeek.ErrInvalidOption)
	assert.Nil(t, result)
}
