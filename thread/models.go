package thread

// ThreadDetail is a forum thread returned by the thread endpoint, including its articles (posts).
type ThreadDetail struct {
	ID          int       `xml:"id,attr"`
	NumArticles int       `xml:"numarticles,attr"`
	Link        string    `xml:"link,attr"`
	Subject     string    `xml:"subject"`
	Articles    []Article `xml:"articles>article"`
}

// Article is a single post within a forum thread.
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
