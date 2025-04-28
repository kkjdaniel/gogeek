package family

type Family struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Type        string `xml:"type,attr"`
	ID          int    `xml:"id,attr"`
	Thumbnail   string `xml:"thumbnail"`
	Image       string `xml:"image"`
	Name        Name   `xml:"name"`
	Description string `xml:"description"`
	Links       []Link `xml:"link"`
}

type Name struct {
	Type      string `xml:"type,attr"`
	SortIndex int    `xml:"sortindex,attr"`
	Value     string `xml:"value,attr"`
}

type Link struct {
	Type    string `xml:"type,attr"`
	ID      int    `xml:"id,attr"`
	Value   string `xml:"value,attr"`
	Inbound bool   `xml:"inbound,attr,omitempty"`
}
