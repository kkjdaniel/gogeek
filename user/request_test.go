package user

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

const mockDataFileValid = "testdata/valid_user_response.xml"

func TestQueryUser(t *testing.T) {
	defer testutils.ActivateMocks()()

	mockURL := constants.UserEndpoint + "?buddies=1&guilds=1&name=johndoe&top=1"
	testutils.SetupMockResponder(t, mockURL, mockDataFileValid)

	client := gogeek.NewClient(gogeek.APIKey("test-key"))
	user, err := Query(context.Background(), client, "johndoe")
	require.NoError(t, err, "Query should not return an error")
	require.NotNil(t, user, "User should not be nil")

	expected := &User{
		ID:               12345,
		Name:             "John Doe",
		FirstName:        ValueField{Value: "John"},
		LastName:         ValueField{Value: "Doe"},
		AvatarLink:       ValueField{Value: "https://example.com/avatars/avatar_default.png"},
		YearRegistered:   IntValueField{Value: 2003},
		LastLogin:        ValueField{Value: "2025-04-04"},
		StateOrProvince:  ValueField{Value: "Example State"},
		Country:          ValueField{Value: "Example Country"},
		WebAddress:       ValueField{Value: "https://example.com/blog"},
		XboxAccount:      ValueField{Value: ""},
		WiiAccount:       ValueField{Value: ""},
		PSNAccount:       ValueField{Value: ""},
		BattleNetAccount: ValueField{Value: ""},
		SteamAccount:     ValueField{Value: ""},
		TradeRating:      IntValueField{Value: 0},
		Buddies: Buddies{
			Total: 5,
			Page:  1,
			Buddy: []Buddy{
				{ID: 1001, Name: "buddy_one"},
				{ID: 1002, Name: "buddy_two"},
				{ID: 1003, Name: "buddy_three"},
				{ID: 1004, Name: "buddy_four"},
				{ID: 1005, Name: "buddy_five"},
			},
		},
		Guilds: Guilds{
			Total: 2,
			Page:  1,
			Guild: []Guild{
				{ID: 2001, Name: "Example Guild One"},
				{ID: 2002, Name: "Example Guild Two"},
			},
		},
		Top: Top{
			Domain: "boardgame",
			Items: []TopItem{
				{Rank: 1, Type: "thing", ID: 3001, Name: "Example Game One"},
				{Rank: 2, Type: "thing", ID: 3002, Name: "Example Game Two"},
			},
		},
	}

	if diff := cmp.Diff(expected, user); diff != "" {
		t.Errorf("User mismatch (-want +got):\n%s", diff)
	}
}

func TestQuery_Error(t *testing.T) {
	testURL := constants.UserEndpoint + "?buddies=1&guilds=1&name=johndoe&top=1"

	queryWrapper := func(url string) (*User, error) {
		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		return Query(context.Background(), client, "johndoe")
	}

	testutils.TestRequestError(t, testURL, queryWrapper)
}

func TestUserOptions(t *testing.T) {
	defaults := func() url.Values {
		params := url.Values{}
		params.Set("buddies", "1")
		params.Set("guilds", "1")
		params.Set("top", "1")
		return params
	}

	t.Run("WithoutBuddies", func(t *testing.T) {
		params := defaults()
		require.NoError(t, WithoutBuddies()(params))
		assert.Empty(t, params.Get("buddies"))
		assert.Equal(t, "1", params.Get("guilds"))
	})

	t.Run("WithoutGuilds", func(t *testing.T) {
		params := defaults()
		require.NoError(t, WithoutGuilds()(params))
		assert.Empty(t, params.Get("guilds"))
	})

	t.Run("WithoutTop", func(t *testing.T) {
		params := defaults()
		require.NoError(t, WithoutTop()(params))
		assert.Empty(t, params.Get("top"))
	})

	t.Run("WithHot", func(t *testing.T) {
		params := defaults()
		require.NoError(t, WithHot()(params))
		assert.Equal(t, "1", params.Get("hot"))
	})

	t.Run("WithDomain", func(t *testing.T) {
		params := defaults()
		require.NoError(t, WithDomain("rpg")(params))
		assert.Equal(t, "rpg", params.Get("domain"))
	})

	t.Run("WithPage", func(t *testing.T) {
		params := defaults()
		require.NoError(t, WithPage(2)(params))
		assert.Equal(t, "2", params.Get("page"))
	})

	t.Run("Invalid values", func(t *testing.T) {
		params := defaults()
		assert.ErrorIs(t, WithDomain("bogus")(params), gogeek.ErrInvalidOption)
		assert.ErrorIs(t, WithPage(0)(params), gogeek.ErrInvalidOption)

		client := gogeek.NewClient(gogeek.APIKey("test-key"))
		result, err := Query(context.Background(), client, "johndoe", WithDomain("bogus"))
		assert.ErrorIs(t, err, gogeek.ErrInvalidOption)
		assert.Nil(t, result)
	})
}
