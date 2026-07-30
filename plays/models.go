package plays

// Plays is a page of a user's logged plays as returned by the plays endpoint.
type Plays struct {
	UserID   int    `xml:"userid,attr"`
	Username string `xml:"username,attr"`
	Total    int    `xml:"total,attr"`
	Page     int    `xml:"page,attr"`
	Plays    []Play `xml:"play"`
}

// Play is a single logged play, including the item played and its players.
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

// Player is a participant in a logged play, with their score and result.
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

// PlayItem is the game or item that a play was logged against.
type PlayItem struct {
	Name       string    `xml:"name,attr"`
	ObjectType string    `xml:"objecttype,attr"`
	ObjectID   int       `xml:"objectid,attr"`
	Subtypes   []Subtype `xml:"subtypes>subtype"`
}

// Subtype holds an XML value attribute naming a played item's subtype (e.g. boardgame).
type Subtype struct {
	Value string `xml:"value,attr"`
}
