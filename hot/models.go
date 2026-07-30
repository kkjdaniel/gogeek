package hot

// HotItems is the list of trending items returned by the hot endpoint.
type HotItems struct {
	Items []HotItem `xml:"item"`
}

// HotItem is a single entry in the BGG hotness ranking.
type HotItem struct {
	ID            int         `xml:"id,attr"`
	Rank          int         `xml:"rank,attr"`
	Name          ValueString `xml:"name"`
	YearPublished ValueInt    `xml:"yearpublished"`
	Thumbnail     ValueString `xml:"thumbnail"`
}

// ValueString holds a string XML value attribute.
type ValueString struct {
	Value string `xml:"value,attr"`
}

// ValueInt holds an integer XML value attribute.
type ValueInt struct {
	Value int `xml:"value,attr"`
}
