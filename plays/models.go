package plays

type Plays struct {
	UserID   int    `xml:"userid,attr"`
	Username string `xml:"username,attr"`
	Total    int    `xml:"total,attr"`
	Page     int    `xml:"page,attr"`
	Plays    []Play `xml:"play"`
}

type Play struct {
	ID         int      `xml:"id,attr"`
	Date       string   `xml:"date,attr"`
	Quantity   int      `xml:"quantity,attr"`
	Length     int      `xml:"length,attr"`
	Incomplete int      `xml:"incomplete,attr"`
	NoWinStats int      `xml:"nowinstats,attr"`
	Location   string   `xml:"location,attr"`
	Item       PlayItem `xml:"item"`
	Comments   string   `xml:"comments"`
	Players    []Player `xml:"players>player"`
}

type Player struct {
	Username      string `xml:"username,attr"`
	UserID        int    `xml:"userid,attr"`
	Name          string `xml:"name,attr"`
	StartPosition string `xml:"startposition,attr"`
	Color         string `xml:"color,attr"`
	Score         int    `xml:"score,attr"`
	New           int    `xml:"new,attr"`
	Rating        int    `xml:"rating,attr"`
	Win           int    `xml:"win,attr"`
}

type PlayItem struct {
	Name       string    `xml:"name,attr"`
	ObjectType string    `xml:"objecttype,attr"`
	ObjectID   int       `xml:"objectid,attr"`
	Subtypes   []Subtype `xml:"subtypes>subtype"`
}

type Subtype struct {
	Value string `xml:"value,attr"`
}
