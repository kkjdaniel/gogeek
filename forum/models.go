package forum

type Forum struct {
	ID           int      `xml:"id,attr"`
	Title        string   `xml:"title,attr"`
	NumThreads   int      `xml:"numthreads,attr"`
	NumPosts     int      `xml:"numposts,attr"`
	LastPostDate string   `xml:"lastpostdate,attr"`
	NoPosting    int      `xml:"noposting,attr"`
	Threads      []Thread `xml:"threads>thread"`
}

type Thread struct {
	ID           int    `xml:"id,attr"`
	Subject      string `xml:"subject,attr"`
	Author       string `xml:"author,attr"`
	NumArticles  int    `xml:"numarticles,attr"`
	PostDate     string `xml:"postdate,attr"`
	LastPostDate string `xml:"lastpostdate,attr"`
}
