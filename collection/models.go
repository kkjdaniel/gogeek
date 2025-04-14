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
	Status        ItemStatus `xml:"status"`
	NumPlays      int        `xml:"numplays"`
	Comment       string     `xml:"comment"`
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
