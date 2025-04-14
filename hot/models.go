package hot

type HotItems struct {
	Items []HotItem `xml:"item"`
}

type HotItem struct {
	ID            int         `xml:"id,attr"`
	Rank          int         `xml:"rank,attr"`
	Name          ValueString `xml:"name"`
	YearPublished ValueInt    `xml:"yearpublished"`
	Thumbnail     ValueString `xml:"thumbnail"`
}

type ValueString struct {
	Value string `xml:"value,attr"`
}

type ValueInt struct {
	Value int `xml:"value,attr"`
}
