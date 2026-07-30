package family

// Family is the response from the BGG family endpoint, containing the
// requested family items.
type Family struct {
	Items []Item `xml:"item"`
}

// Item is a single BGG family, including its name, description, and links to
// the things (e.g. board games) that belong to it.
type Item struct {
	Type        string `xml:"type,attr"`
	ID          int    `xml:"id,attr"`
	Thumbnail   string `xml:"thumbnail"`
	Image       string `xml:"image"`
	Name        Name   `xml:"name"`
	Description string `xml:"description"`
	Links       []Link `xml:"link"`
}

// Name is a family's name, held in an XML value attribute.
type Name struct {
	Type      string `xml:"type,attr"`
	SortIndex int    `xml:"sortindex,attr"`
	Value     string `xml:"value,attr"`
}

// Link is a reference from a family to a related BGG item, such as a board
// game within the family.
type Link struct {
	Type    string `xml:"type,attr"`
	ID      int    `xml:"id,attr"`
	Value   string `xml:"value,attr"`
	Inbound bool   `xml:"inbound,attr,omitempty"`
}
