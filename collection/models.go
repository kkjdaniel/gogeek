package collection

type Collection struct {
	TotalItems int              `xml:"totalitems,attr"`
	PubDate    string           `xml:"pubdate,attr"`
	Items      []CollectionItem `xml:"item"`
}

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

type ItemStats struct {
	MinPlayers  int          `xml:"minplayers,attr"`
	MaxPlayers  int          `xml:"maxplayers,attr"`
	MinPlayTime int          `xml:"minplaytime,attr"`
	MaxPlayTime int          `xml:"maxplaytime,attr"`
	PlayingTime int          `xml:"playingtime,attr"`
	NumOwned    int          `xml:"numowned,attr"`
	Rating      StatsRating  `xml:"rating"`
}

type StatsRating struct {
	Value        string      `xml:"value,attr"`
	UsersRated   RatingValue `xml:"usersrated"`
	Average      RatingValue `xml:"average"`
	BayesAverage RatingValue `xml:"bayesaverage"`
	StdDev       RatingValue `xml:"stddev"`
	Median       RatingValue `xml:"median"`
	Ranks        []StatsRank `xml:"ranks>rank"`
}

type RatingValue struct {
	Value string `xml:"value,attr"`
}

type StatsRank struct {
	Type         string `xml:"type,attr"`
	ID           int    `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	FriendlyName string `xml:"friendlyname,attr"`
	Value        string `xml:"value,attr"`
	BayesAverage string `xml:"bayesaverage,attr"`
}

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
