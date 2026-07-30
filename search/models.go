package search

// SearchResults is the set of items returned by the search endpoint.
type SearchResults struct {
	Total int            `xml:"total,attr"`
	Items []SearchResult `xml:"item"`
}

// SearchResult is a single item matched by a BGG search query.
type SearchResult struct {
	ID            int              `xml:"id,attr"`
	Type          string           `xml:"type,attr"`
	Name          Name             `xml:"name"`
	YearPublished YearPublishedTag `xml:"yearpublished"`
}

// Name holds an item's name as an XML value attribute, with its type (primary or alternate).
type Name struct {
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

// YearPublishedTag holds an item's publication year as an XML value attribute.
type YearPublishedTag struct {
	Value int `xml:"value,attr"`
}
