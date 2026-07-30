package forumlist

// ForumList is the list of forums for a BGG thing or family, as returned by
// the forumlist endpoint.
type ForumList struct {
	Type   string  `xml:"type,attr"`
	ID     int     `xml:"id,attr"`
	Forums []Forum `xml:"forum"`
}

// Forum is a single forum entry in a forum list, with its title, description,
// and thread and post counts.
type Forum struct {
	ID           int    `xml:"id,attr"`
	GroupID      int    `xml:"groupid,attr"`
	Title        string `xml:"title,attr"`
	NoPosting    int    `xml:"noposting,attr"`
	Description  string `xml:"description,attr"`
	NumThreads   int    `xml:"numthreads,attr"`
	NumPosts     int    `xml:"numposts,attr"`
	LastPostDate string `xml:"lastpostdate,attr"`
}
