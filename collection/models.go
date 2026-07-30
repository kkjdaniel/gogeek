package collection

// Collection is a user's BGG collection as returned by the collection endpoint.
type Collection struct {
	TotalItems int              `xml:"totalitems,attr"`
	PubDate    string           `xml:"pubdate,attr"`
	Items      []CollectionItem `xml:"item"`
}

// CollectionItem is a single item (e.g. a board game) in a user's BGG collection.
type CollectionItem struct {
	ObjectType    string     `xml:"objecttype,attr"`
	ObjectID      int        `xml:"objectid,attr"`
	Subtype       string     `xml:"subtype,attr"`
	CollectionID  int        `xml:"collid,attr"`
	Name          string     `xml:"name"`
	YearPublished int        `xml:"yearpublished"`
	Image         string     `xml:"image"`
	Thumbnail     string     `xml:"thumbnail"`
	Stats         *ItemStats `xml:"stats"`
	Status        ItemStatus `xml:"status"`
	NumPlays      int        `xml:"numplays"`
	Comment       string     `xml:"comment"`
}

// ItemStats holds gameplay statistics for a collection item, such as player
// counts, play times, and community rating data.
type ItemStats struct {
	MinPlayers  int         `xml:"minplayers,attr"`
	MaxPlayers  int         `xml:"maxplayers,attr"`
	MinPlayTime int         `xml:"minplaytime,attr"`
	MaxPlayTime int         `xml:"maxplaytime,attr"`
	PlayingTime int         `xml:"playingtime,attr"`
	NumOwned    int         `xml:"numowned,attr"`
	Rating      StatsRating `xml:"rating"`
}

// StatsRating holds the user's rating for an item along with aggregate
// community rating statistics and rank information.
type StatsRating struct {
	Value        string      `xml:"value,attr"`
	UsersRated   RatingValue `xml:"usersrated"`
	Average      RatingValue `xml:"average"`
	BayesAverage RatingValue `xml:"bayesaverage"`
	StdDev       RatingValue `xml:"stddev"`
	Median       RatingValue `xml:"median"`
	Ranks        []StatsRank `xml:"ranks>rank"`
}

// RatingValue holds a rating statistic exposed as an XML value attribute.
type RatingValue struct {
	Value string `xml:"value,attr"`
}

// StatsRank is an item's rank within a BGG ranking category, such as the
// overall board game rank or a subdomain rank.
type StatsRank struct {
	Type         string `xml:"type,attr"`
	ID           int    `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	FriendlyName string `xml:"friendlyname,attr"`
	Value        string `xml:"value,attr"`
	BayesAverage string `xml:"bayesaverage,attr"`
}

// ItemStatus holds the ownership and wishlist status flags for a collection
// item, such as owned, for trade, want to play, and preordered.
type ItemStatus struct {
	Own          int    `xml:"own,attr"`
	PrevOwned    int    `xml:"prevowned,attr"`
	ForTrade     int    `xml:"fortrade,attr"`
	Want         int    `xml:"want,attr"`
	WantToPlay   int    `xml:"wanttoplay,attr"`
	WantToBuy    int    `xml:"wanttobuy,attr"`
	Wishlist     int    `xml:"wishlist,attr"`
	Preordered   int    `xml:"preordered,attr"`
	LastModified string `xml:"lastmodified,attr"`
}
