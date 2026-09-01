//go:build contract

package contract

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	gogeek "github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/collection"
	"github.com/kkjdaniel/gogeek/v3/constants"
	"github.com/kkjdaniel/gogeek/v3/family"
	"github.com/kkjdaniel/gogeek/v3/forum"
	"github.com/kkjdaniel/gogeek/v3/forumlist"
	"github.com/kkjdaniel/gogeek/v3/guild"
	"github.com/kkjdaniel/gogeek/v3/hot"
	"github.com/kkjdaniel/gogeek/v3/plays"
	"github.com/kkjdaniel/gogeek/v3/search"
	"github.com/kkjdaniel/gogeek/v3/thing"
	"github.com/kkjdaniel/gogeek/v3/thread"
	"github.com/kkjdaniel/gogeek/v3/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Well-known stable IDs for contract testing.
const (
	catanThingID        = 13          // Catan - a well-established game that won't be removed
	catanFamilyID       = 3           // Catan family
	carcassonneFamilyID = 2           // Carcassonne family
	knownUsername       = "kkjdaniel" // An established BGG user
	knownGuildID        = 1           // First guild on BGG
	knownForumID        = 19          // A well-known BGG forum
	knownThreadID       = 100000      // A long-standing thread
	knownForumObjID     = 13          // Catan's forum list (thing type)
	unrankedThingID     = 399366      // An obscure unranked item
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
		return gogeek.NewClient(gogeek.APIKey(apiKey))
	}
	if cookie := os.Getenv("BGG_COOKIE"); cookie != "" {
		return gogeek.NewClient(gogeek.Cookie(cookie))
	}

	t.Fatal("BGG_API_KEY or BGG_COOKIE environment variable must be set to run contract tests")
	return nil
}

func mustParseDate(t *testing.T, date string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", date)
	require.NoError(t, err)
	return parsed
}

func TestContract_Thing(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&stats=1", constants.ThingEndpoint, catanThingID)

	result, err := thing.Query(ctx, client, []int{catanThingID})
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
	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing")
}

func TestContract_Thing_WithVideos(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&stats=1&videos=1", constants.ThingEndpoint, catanThingID)

	result, err := thing.Query(ctx, client, []int{catanThingID}, thing.WithVideos())
	require.NoError(t, err, "thing.Query with videos should not error")
	require.NotEmpty(t, result.Items, "should return at least one item")

	item := result.Items[0]
	require.NotNil(t, item.Videos, "a game as established as Catan should carry videos")
	require.NotEmpty(t, item.Videos.Videos, "the videos element should hold entries")
	assert.Greater(t, item.Videos.Total, 0, "total should be positive")

	video := item.Videos.Videos[0]
	assert.Greater(t, video.ID, 0, "video should have an ID")
	assert.NotEmpty(t, video.Title, "video should have a title")
	assert.NotEmpty(t, video.Category, "video should have a category")
	assert.NotEmpty(t, video.Link, "video should have a link")
	assert.NotEmpty(t, video.Username, "video should have an uploader")
	assert.NotEmpty(t, video.PostDate, "video should have a post date")

	// Field coverage: check the API hasn't added fields our model doesn't capture
	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing?videos=1")
}

func TestContract_Thing_WithVersions(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&stats=1&versions=1", constants.ThingEndpoint, catanThingID)

	result, err := thing.Query(ctx, client, []int{catanThingID}, thing.WithVersions())
	require.NoError(t, err, "thing.Query with versions should not error")
	require.NotEmpty(t, result.Items, "should return at least one item")

	item := result.Items[0]
	require.NotEmpty(t, item.Versions, "a game as established as Catan should carry versions")

	version := item.Versions[0]
	assert.Greater(t, version.ID, 0, "version should have an ID")
	assert.NotEmpty(t, version.Type, "version should have a type")
	assert.NotEmpty(t, version.Name, "version should have a name")
	assert.NotEmpty(t, version.Links, "version should carry links")

	// Field coverage: check the API hasn't added fields our model doesn't capture
	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing?versions=1")
}

func TestContract_Thing_WithComments(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&comments=1&stats=1", constants.ThingEndpoint, catanThingID)

	result, err := thing.Query(ctx, client, []int{catanThingID}, thing.WithComments())
	require.NoError(t, err, "thing.Query with comments should not error")
	require.NotEmpty(t, result.Items, "should return at least one item")

	item := result.Items[0]
	require.NotNil(t, item.Comments, "a game as established as Catan should carry comments")
	require.NotEmpty(t, item.Comments.Comments, "the comments element should hold entries")
	assert.Equal(t, 1, item.Comments.Page, "the first page should be returned by default")
	assert.Greater(t, item.Comments.Total, 0, "totalitems should be positive")
	assert.NotEmpty(t, item.Comments.Comments[0].Username, "comment should have a username")

	// Field coverage: check the API hasn't added fields our model doesn't capture
	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing?comments=1")
}

func TestContract_Thing_WithRatingComments(t *testing.T) {
	client := newClient(t)

	result, err := thing.Query(context.Background(), client, []int{catanThingID}, thing.WithRatingComments(), thing.WithPageSize(10))
	require.NoError(t, err, "thing.Query with rating comments should not error")
	require.NotEmpty(t, result.Items, "should return at least one item")

	item := result.Items[0]
	require.NotNil(t, item.Comments, "ratings arrive in the comments node")
	require.NotEmpty(t, item.Comments.Comments, "the comments element should hold entries")
	assert.NotEqual(t, "N/A", item.Comments.Comments[0].Rating, "rating comments should all carry a rating")
}

func TestContract_Thing_WithMarketplace(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&marketplace=1&stats=1", constants.ThingEndpoint, catanThingID)

	result, err := thing.Query(ctx, client, []int{catanThingID}, thing.WithMarketplace())
	require.NoError(t, err, "thing.Query with marketplace should not error")
	require.NotEmpty(t, result.Items, "should return at least one item")

	item := result.Items[0]
	require.NotEmpty(t, item.Marketplace, "a game as established as Catan should carry listings")

	listing := item.Marketplace[0]
	assert.NotEmpty(t, listing.Price.Currency, "listing should have a currency")
	assert.NotEmpty(t, listing.Condition.Value, "listing should have a condition")
	assert.NotEmpty(t, listing.Link.Href, "listing should link to the marketplace page")

	// Field coverage: check the API hasn't added fields our model doesn't capture
	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing?marketplace=1")
}

func TestContract_Thing_WithType(t *testing.T) {
	client := newClient(t)

	// Catan is a boardgame, so filtering for expansions must exclude it.
	result, err := thing.Query(context.Background(), client, []int{catanThingID}, thing.WithType("boardgameexpansion"))
	require.NoError(t, err, "thing.Query with a type filter should not error")
	assert.Empty(t, result.Items, "a boardgame should be filtered out by an expansion-only query")
}

func TestContract_Thing_MultipleIDs(t *testing.T) {
	client := newClient(t)

	ids := []int{catanThingID, 822} // Catan and Carcassonne
	result, err := thing.Query(context.Background(), client, ids)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Items), 2, "should return multiple items")
}

func TestContract_Thing_Unranked(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&stats=1", constants.ThingEndpoint, unrankedThingID)

	result, err := thing.Query(ctx, client, []int{unrankedThingID})
	require.NoError(t, err, "thing.Query should not error for unranked item")
	require.NotNil(t, result)
	require.NotEmpty(t, result.Items, "should return the item")

	item := result.Items[0]
	assert.Equal(t, unrankedThingID, item.ID)
	assert.NotEmpty(t, item.Name, "unranked item should have a name")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thing.Items{}, "thing-unranked")
}

func TestContract_Thing_AutoChunk(t *testing.T) {
	client := newClient(t)

	// 21 well-known IDs force two rate-limited batches (20 + 1).
	ids := []int{13, 822, 68448, 30549, 36218, 84876, 12333, 3076, 31260,
		2651, 9209, 178900, 167791, 174430, 220308, 169786, 120677, 102794,
		28720, 25613, 1406}
	result, err := thing.Query(context.Background(), client, ids)
	require.NoError(t, err, "thing.Query should transparently chunk 21 IDs")
	require.NotNil(t, result)
	assert.Equal(t, len(ids), len(result.Items), "all chunked results should be merged")
}

func TestContract_Search(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := constants.SearchEndpoint + "?query=Catan"

	result, err := search.Query(ctx, client, "Catan")
	require.NoError(t, err, "search.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	assert.Greater(t, result.Total, 0, "should find results for 'Catan'")
	require.NotEmpty(t, result.Items, "items slice should not be empty")

	item := result.Items[0]
	assert.Greater(t, item.ID, 0, "search result should have a positive ID")
	assert.NotEmpty(t, item.Name.Value, "search result should have a name")
	assert.NotEmpty(t, item.Type, "search result should have a type")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, search.SearchResults{}, "search")
}

func TestContract_Search_Exact(t *testing.T) {
	client := newClient(t)

	result, err := search.Query(context.Background(), client, "Catan", search.WithExact())
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.Total, 0, "exact search for 'Catan' should return results")
}

func TestContract_Search_Type(t *testing.T) {
	client := newClient(t)

	result, err := search.Query(context.Background(), client, "Catan",
		search.WithType("boardgameexpansion"))
	require.NoError(t, err, "search.Query with type filter should not error")
	require.NotNil(t, result)
	require.NotEmpty(t, result.Items, "type-filtered search should return results")
	for _, item := range result.Items {
		assert.Equal(t, "boardgameexpansion", item.Type,
			"type filter should restrict results to the requested type")
	}
}

func TestContract_Hot(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := constants.HotEndpoint + "?type=boardgame"

	result, err := hot.Query(ctx, client, hot.ItemTypeBoardGame)
	require.NoError(t, err, "hot.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	require.NotEmpty(t, result.Items, "hot items should not be empty")

	item := result.Items[0]
	assert.Greater(t, item.Rank, 0, "hot item should have a positive rank")
	assert.NotEmpty(t, item.Name.Value, "hot item should have a name")
	assert.NotEmpty(t, item.Thumbnail.Value, "hot item should have a thumbnail")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, hot.HotItems{}, "hot")
}

func TestContract_User(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := constants.UserEndpoint + "?name=" + knownUsername

	result, err := user.Query(ctx, client, knownUsername)
	require.NoError(t, err, "user.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "user should have a positive ID")
	assert.Equal(t, knownUsername, result.Name, "username should match")
	assert.Greater(t, result.YearRegistered.Value, 0, "user should have a registration year")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, user.User{}, "user")
}

func TestContract_User_Options(t *testing.T) {
	client := newClient(t)

	result, err := user.Query(context.Background(), client, knownUsername,
		user.WithoutBuddies(),
		user.WithoutGuilds(),
		user.WithHot())
	require.NoError(t, err, "user.Query with options should not error")
	require.NotNil(t, result)

	assert.Equal(t, knownUsername, result.Name)
	assert.Empty(t, result.Buddies.Buddy, "buddies should be omitted when opted out")
	assert.Empty(t, result.Guilds.Guild, "guilds should be omitted when opted out")
}

func TestContract_Collection(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := constants.CollectionEndpoint + "?username=" + knownUsername + "&stats=1"

	result, err := collection.Query(ctx, client, knownUsername, collection.WithStats())
	require.NoError(t, err, "collection.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	assert.Greater(t, result.TotalItems, 0, "collection should have items")
	require.NotEmpty(t, result.Items, "items slice should not be empty")

	item := result.Items[0]
	assert.Greater(t, item.ObjectID, 0, "collection item should have an object ID")
	assert.NotEmpty(t, item.Name, "collection item should have a name")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, collection.Collection{}, "collection")
}

func TestContract_Family(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&type=boardgamefamily", constants.FamilyEndpoint, catanFamilyID)

	result, err := family.Query(ctx, client, catanFamilyID, "boardgamefamily")
	require.NoError(t, err, "family.Query should not error")
	require.NotNil(t, result, "result should not be nil")
	require.NotEmpty(t, result.Items, "family should have items")

	item := result.Items[0]
	assert.Greater(t, item.ID, 0, "family item should have a positive ID")
	assert.NotEmpty(t, item.Name.Value, "family item should have a name")
	assert.NotEmpty(t, item.Type, "family item should have a type")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, family.Family{}, "family")
}

func TestContract_Family_MultipleIDs(t *testing.T) {
	client := newClient(t)

	// Catan and Carcassonne families in one request.
	result, err := family.Query(context.Background(), client, catanFamilyID, "boardgamefamily", carcassonneFamilyID)
	require.NoError(t, err, "family.Query with multiple IDs should not error")
	require.Len(t, result.Items, 2, "both families should be returned")
}

func TestContract_Forum(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d", constants.ForumEndpoint, knownForumID)

	result, err := forum.Query(ctx, client, knownForumID)
	require.NoError(t, err, "forum.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "forum should have a positive ID")
	assert.NotEmpty(t, result.Title, "forum should have a title")
	assert.Greater(t, result.NumThreads, 0, "forum should have threads")
	assert.Greater(t, result.NumPosts, 0, "forum should have posts")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, forum.Forum{}, "forum")
}

func TestContract_ForumList(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&type=thing", constants.ForumListEndpoint, knownForumObjID)

	result, err := forumlist.Query(ctx, client, knownForumObjID, "thing")
	require.NoError(t, err, "forumlist.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "forumlist should have a positive ID")
	assert.NotEmpty(t, result.Type, "forumlist should have a type")
	assert.NotEmpty(t, result.Forums, "forumlist should have forums")

	f := result.Forums[0]
	assert.Greater(t, f.ID, 0, "forum in list should have a positive ID")
	assert.NotEmpty(t, f.Title, "forum in list should have a title")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, forumlist.ForumList{}, "forumlist")
}

func TestContract_Guild(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d", constants.GuildEndpoint, knownGuildID)

	result, err := guild.Query(ctx, client, knownGuildID)
	require.NoError(t, err, "guild.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "guild should have a positive ID")
	assert.NotEmpty(t, result.Name, "guild should have a name")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, guild.Guild{}, "guild")
}

func TestContract_Guild_Members(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d&members=1", constants.GuildEndpoint, knownGuildID)

	result, err := guild.Query(ctx, client, knownGuildID,
		guild.WithMembers(),
		guild.WithSort("username"))
	require.NoError(t, err, "guild.Query with members should not error")
	require.NotNil(t, result)
	require.NotNil(t, result.Members, "member roster should be present when requested")
	assert.Greater(t, result.Members.Count, 0, "guild should have members")
	assert.NotEmpty(t, result.Members.Members, "member list should not be empty")
	assert.NotEmpty(t, result.Members.Members[0].Name, "members should have usernames")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, guild.Guild{}, "guild-members")
}

func TestContract_Plays(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := constants.PlaysEndpoint + "?username=" + knownUsername

	result, err := plays.Query(ctx, client, knownUsername)
	require.NoError(t, err, "plays.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.UserID, 0, "plays should have a positive user ID")
	assert.Equal(t, knownUsername, result.Username, "username should match")
	assert.GreaterOrEqual(t, result.Total, 0, "total plays should be non-negative")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, plays.Plays{}, "plays")
}

func TestContract_Plays_Options(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	t.Run("page", func(t *testing.T) {
		result, err := plays.Query(ctx, client, knownUsername, plays.WithPage(2))
		require.NoError(t, err, "plays.Query with page should not error")
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Page, "BGG should echo the requested page")
	})

	t.Run("id and subtype", func(t *testing.T) {
		result, err := plays.Query(ctx, client, knownUsername,
			plays.WithID(catanThingID),
			plays.WithType("thing"),
			plays.WithSubtype("boardgame"))
		require.NoError(t, err, "plays.Query with id/type/subtype should not error")
		require.NotNil(t, result)
		for _, p := range result.Plays {
			assert.Equal(t, catanThingID, p.Item.ObjectID,
				"id filter should restrict plays to the requested item")
		}
	})

	t.Run("date range", func(t *testing.T) {
		minDate := "2024-01-01"
		maxDate := "2024-12-31"
		result, err := plays.Query(ctx, client, knownUsername,
			plays.WithMinDate(mustParseDate(t, minDate)),
			plays.WithMaxDate(mustParseDate(t, maxDate)))
		require.NoError(t, err, "plays.Query with date range should not error")
		require.NotNil(t, result)
		for _, p := range result.Plays {
			assert.GreaterOrEqual(t, p.Date, minDate, "plays should be on or after mindate")
			assert.LessOrEqual(t, p.Date, maxDate, "plays should be on or before maxdate")
		}
	})
}

func TestContract_Thread(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	url := fmt.Sprintf("%s?id=%d", constants.ThreadEndpoint, knownThreadID)

	result, err := thread.Query(ctx, client, knownThreadID)
	require.NoError(t, err, "thread.Query should not error")
	require.NotNil(t, result, "result should not be nil")

	assert.Greater(t, result.ID, 0, "thread should have a positive ID")
	assert.NotEmpty(t, result.Subject, "thread should have a subject")
	assert.Greater(t, result.NumArticles, 0, "thread should have articles")
	require.NotEmpty(t, result.Articles, "articles slice should not be empty")

	rawXML, err := fetchRawXML(ctx, client, url)
	require.NoError(t, err, "fetching raw XML for coverage check")
	assertFieldCoverage(t, rawXML, thread.ThreadDetail{}, "thread")
}

func TestContract_Thread_Count(t *testing.T) {
	client := newClient(t)

	result, err := thread.Query(context.Background(), client, knownThreadID,
		thread.WithCount(1))
	require.NoError(t, err, "thread.Query with count should not error")
	require.NotNil(t, result)
	assert.LessOrEqual(t, len(result.Articles), 1, "count should limit returned articles")
}
