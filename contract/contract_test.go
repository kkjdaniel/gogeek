//go:build contract

package contract

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	gogeek "github.com/kkjdaniel/gogeek/v2"
	"github.com/kkjdaniel/gogeek/v2/collection"
	"github.com/kkjdaniel/gogeek/v2/constants"
	"github.com/kkjdaniel/gogeek/v2/family"
	"github.com/kkjdaniel/gogeek/v2/forum"
	"github.com/kkjdaniel/gogeek/v2/forumlist"
	"github.com/kkjdaniel/gogeek/v2/guild"
	"github.com/kkjdaniel/gogeek/v2/hot"
	"github.com/kkjdaniel/gogeek/v2/plays"
	"github.com/kkjdaniel/gogeek/v2/search"
	"github.com/kkjdaniel/gogeek/v2/thing"
	"github.com/kkjdaniel/gogeek/v2/thread"
	"github.com/kkjdaniel/gogeek/v2/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Well-known stable IDs for contract testing.
const (
	catanThingID    = 13         // Catan - a well-established game that won't be removed
	catanFamilyID   = 3          // Catan family
	knownUsername   = "kkjdaniel" // An established BGG user
	knownGuildID    = 1          // First guild on BGG
	knownForumID    = 19         // A well-known BGG forum
	knownThreadID   = 100000     // A long-standing thread
	knownForumObjID = 13         // Catan's forum list (thing type)
	unrankedThingID = 399366     // An obscure unranked item
)

func TestMain(m *testing.M) {
	// Load .env from project root (one level up from contract/)
	_, filename, _, _ := runtime.Caller(0)
	envPath := filepath.Join(filepath.Dir(filename), "..", ".env")
	if f, err := os.Open(envPath); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if k, v, ok := strings.Cut(line, "="); ok {
				if os.Getenv(k) == "" {
					os.Setenv(k, v)
				}
			}
		}
		f.Close()
	}
	os.Exit(m.Run())
}

func newClient(t *testing.T) *gogeek.Client {
	t.Helper()

	if apiKey := os.Getenv("BGG_API_KEY"); apiKey != "" {
		return gogeek.NewClient(gogeek.WithAPIKey(apiKey))
	}
	if cookie := os.Getenv("BGG_COOKIE"); cookie != "" {
		return gogeek.NewClient(gogeek.WithCookie(cookie))
	}

	t.Fatal("BGG_API_KEY or BGG_COOKIE environment variable must be set to run contract tests")
	return nil
}

func TestContract_Thing(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d&stats=1", constants.ThingEndpoint, catanThingID)

	result, err := thing.Query(client, []int{catanThingID})
	require.NoError(t, err, "thing.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	require.NotEmpty(t, result.Items, "should return at least one item")

	item := result.Items[0]
	assert.Equal(t, catanThingID, item.ID, "item ID should match requested ID")
	assert.NotEmpty(t, item.Name, "item should have names")
	assert.NotEmpty(t, item.Type, "item should have a type")
	assert.NotEmpty(t, item.Description, "item should have a description")
	assert.NotEmpty(t, item.Thumbnail, "item should have a thumbnail URL")
	assert.NotEmpty(t, item.Image, "item should have an image URL")
	assert.Greater(t, item.YearPublished.Value, 0, "year published should be positive")
	assert.Greater(t, item.MinPlayers.Value, 0, "min players should be positive")
	assert.Greater(t, item.MaxPlayers.Value, 0, "max players should be positive")
	assert.Greater(t, item.PlayingTime.Value, 0, "playing time should be positive")
	assert.NotEmpty(t, item.Links, "item should have links")

	found := false
	for _, n := range item.Name {
		if n.Type == "primary" {
			assert.NotEmpty(t, n.Value, "primary name should have a value")
			found = true
		}
	}
	assert.True(t, found, "item should have a primary name")

	// Field coverage: check the API hasn't added fields our model doesn't capture
	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing")
}

func TestContract_Thing_MultipleIDs(t *testing.T) {
	client := newClient(t)

	ids := []int{catanThingID, 822} // Catan and Carcassonne
	result, err := thing.Query(client, ids)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Items), 2, "should return multiple items")
}

func TestContract_Thing_Unranked(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d&stats=1", constants.ThingEndpoint, unrankedThingID)

	result, err := thing.Query(client, []int{unrankedThingID})
	require.NoError(t, err, "thing.Query should not error for unranked item")
	require.NotNil(t, result)
	require.NotEmpty(t, result.Items, "should return the item")

	item := result.Items[0]
	assert.Equal(t, unrankedThingID, item.ID)
	assert.NotEmpty(t, item.Name, "unranked item should have a name")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing-unranked")
}

func TestContract_Search(t *testing.T) {
	client := newClient(t)
	url := constants.SearchEndpoint + "?query=Catan"

	result, err := search.Query(client, "Catan")
	require.NoError(t, err, "search.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	assert.Greater(t, result.Total, 0, "should find results for 'Catan'")
	require.NotEmpty(t, result.Items, "items slice should not be empty")

	item := result.Items[0]
	assert.Greater(t, item.ID, 0, "search result should have a positive ID")
	assert.NotEmpty(t, item.Name.Value, "search result should have a name")
	assert.NotEmpty(t, item.Type, "search result should have a type")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, search.SearchResults{}, "search")
}

func TestContract_Search_Exact(t *testing.T) {
	client := newClient(t)

	result, err := search.Query(client, "Catan", true)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.Total, 0, "exact search for 'Catan' should return results")
}

func TestContract_Hot(t *testing.T) {
	client := newClient(t)
	url := constants.HotEndpoint + "?type=boardgame"

	result, err := hot.Query(client, hot.ItemTypeBoardGame)
	require.NoError(t, err, "hot.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	require.NotEmpty(t, result.Items, "hot items should not be empty")

	item := result.Items[0]
	assert.Greater(t, item.Rank, 0, "hot item should have a positive rank")
	assert.NotEmpty(t, item.Name.Value, "hot item should have a name")
	assert.NotEmpty(t, item.Thumbnail.Value, "hot item should have a thumbnail")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, hot.HotItems{}, "hot")
}

func TestContract_User(t *testing.T) {
	client := newClient(t)
	url := constants.UserEndpoint + "?name=" + knownUsername

	result, err := user.Query(client, knownUsername)
	require.NoError(t, err, "user.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "user should have a positive ID")
	assert.Equal(t, knownUsername, result.Name, "username should match")
	assert.Greater(t, result.YearRegistered.Value, 0, "user should have a registration year")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, user.User{}, "user")
}

func TestContract_Collection(t *testing.T) {
	client := newClient(t)
	url := constants.CollectionEndpoint + "?username=" + knownUsername + "&stats=1"

	result, err := collection.Query(client, knownUsername, collection.WithStats())
	require.NoError(t, err, "collection.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	assert.Greater(t, result.TotalItems, 0, "collection should have items")
	require.NotEmpty(t, result.Items, "items slice should not be empty")

	item := result.Items[0]
	assert.Greater(t, item.ObjectID, 0, "collection item should have an object ID")
	assert.NotEmpty(t, item.Name, "collection item should have a name")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, collection.Collection{}, "collection")
}

func TestContract_Family(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d&type=boardgamefamily", constants.FamilyEndpoint, catanFamilyID)

	result, err := family.Query(client, catanFamilyID, "boardgamefamily")
	require.NoError(t, err, "family.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	require.NotEmpty(t, result.Items, "family should have items")

	item := result.Items[0]
	assert.Greater(t, item.ID, 0, "family item should have a positive ID")
	assert.NotEmpty(t, item.Name.Value, "family item should have a name")
	assert.NotEmpty(t, item.Type, "family item should have a type")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, family.Family{}, "family")
}

func TestContract_Forum(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d", constants.ForumEndpoint, knownForumID)

	result, err := forum.Query(client, knownForumID)
	require.NoError(t, err, "forum.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "forum should have a positive ID")
	assert.NotEmpty(t, result.Title, "forum should have a title")
	assert.Greater(t, result.NumThreads, 0, "forum should have threads")
	assert.Greater(t, result.NumPosts, 0, "forum should have posts")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, forum.Forum{}, "forum")
}

func TestContract_ForumList(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d&type=thing", constants.ForumListEndpoint, knownForumObjID)

	result, err := forumlist.Query(client, knownForumObjID, "thing")
	require.NoError(t, err, "forumlist.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "forumlist should have a positive ID")
	assert.NotEmpty(t, result.Type, "forumlist should have a type")
	assert.NotEmpty(t, result.Forums, "forumlist should have forums")

	f := result.Forums[0]
	assert.Greater(t, f.ID, 0, "forum in list should have a positive ID")
	assert.NotEmpty(t, f.Title, "forum in list should have a title")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, forumlist.ForumList{}, "forumlist")
}

func TestContract_Guild(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d", constants.GuildEndpoint, knownGuildID)

	result, err := guild.Query(client, knownGuildID)
	require.NoError(t, err, "guild.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "guild should have a positive ID")
	assert.NotEmpty(t, result.Name, "guild should have a name")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, guild.Guild{}, "guild")
}

func TestContract_Plays(t *testing.T) {
	client := newClient(t)
	url := constants.PlaysEndpoint + "?username=" + knownUsername

	result, err := plays.Query(client, knownUsername)
	require.NoError(t, err, "plays.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.UserID, 0, "plays should have a positive user ID")
	assert.Equal(t, knownUsername, result.Username, "username should match")
	assert.GreaterOrEqual(t, result.Total, 0, "total plays should be non-negative")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, plays.Plays{}, "plays")
}

func TestContract_Thread(t *testing.T) {
	client := newClient(t)
	url := fmt.Sprintf("%s?id=%d", constants.ThreadEndpoint, knownThreadID)

	result, err := thread.Query(client, knownThreadID)
	require.NoError(t, err, "thread.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "thread should have a positive ID")
	assert.NotEmpty(t, result.Subject, "thread should have a subject")
	assert.Greater(t, result.NumArticles, 0, "thread should have articles")
	require.NotEmpty(t, result.Articles, "articles slice should not be empty")

	rawXML, err := fetchRawXML(client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thread.ThreadDetail{}, "thread")
}
