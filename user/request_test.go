package user

import (
	"testing"

	"github.com/kkjdaniel/gogeek/constants"
	"github.com/kkjdaniel/gogeek/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const mockDataFileValid = "testdata/valid_user_response.xml"

func TestQueryUser(t *testing.T) {
	defer testutils.ActivateMocks()()

	url := constants.UserEndpoint + "?name=johndoe&buddies=1&guilds=1&top=1"
	testutils.SetupMockResponder(t, url, mockDataFileValid)

	user, err := Query("johndoe")
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
