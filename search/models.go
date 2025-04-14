package search

type SearchResults struct {
	Total int            `xml:"total,attr"`
	Items []SearchResult `xml:"item"`
}

type SearchResult struct {
	ID            int              `xml:"id,attr"`
	Type          string           `xml:"type,attr"`
	Name          Name             `xml:"name"`
	YearPublished YearPublishedTag `xml:"yearpublished"`
}

type Name struct {
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type YearPublishedTag struct {
	Value int `xml:"value,attr"`
}
