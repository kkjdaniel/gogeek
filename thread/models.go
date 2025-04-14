package thread

type ThreadDetail struct {
	ID          int       `xml:"id,attr"`
	NumArticles int       `xml:"numarticles,attr"`
	Link        string    `xml:"link,attr"`
	Subject     string    `xml:"subject"`
	Articles    []Article `xml:"articles>article"`
}

type Article struct {
	ID       int    `xml:"id,attr"`
	Username string `xml:"username,attr"`
	Link     string `xml:"link,attr"`
	PostDate string `xml:"postdate,attr"`
	EditDate string `xml:"editdate,attr"`
	NumEdits int    `xml:"numedits,attr"`
	Subject  string `xml:"subject"`
	Body     string `xml:"body"`
}
