package forumlist

type ForumList struct {
	Type   string  `xml:"type,attr"`
	ID     int     `xml:"id,attr"`
	Forums []Forum `xml:"forum"`
}

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
