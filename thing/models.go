package thing

type Items struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Type          string        `xml:"type,attr"`
	ID            int           `xml:"id,attr"`
	Name          []Name        `xml:"name"`
	Description   string        `xml:"description"`
	YearPublished YearPublished `xml:"yearpublished"`
	MinPlayers    Value         `xml:"minplayers"`
	MaxPlayers    Value         `xml:"maxplayers"`
	PlayingTime   Value         `xml:"playingtime"`
	MinPlayTime   Value         `xml:"minplaytime"`
	MaxPlayTime   Value         `xml:"maxplaytime"`
	MinAge        Value         `xml:"minage"`
	Thumbnail     string        `xml:"thumbnail"`
	Image         string        `xml:"image"`
	Links         []Link        `xml:"link"`
	Statistics    *Statistics   `xml:"statistics>ratings"`
	Polls         []Poll        `xml:"poll"`
	PollSummaries []PollSummary `xml:"poll-summary"`
}

type Name struct {
	Type      string `xml:"type,attr"`
	SortIndex int    `xml:"sortindex,attr"`
	Value     string `xml:"value,attr"`
}

type YearPublished struct {
	Value int `xml:"value,attr"`
}

type Value struct {
	Value int `xml:"value,attr"`
}

type Link struct {
	Type  string `xml:"type,attr"`
	ID    int    `xml:"id,attr"`
	Value string `xml:"value,attr"`
}

type Statistics struct {
	UsersRated    Value   `xml:"usersrated"`
	AverageRating float64 `xml:"average"`
	BayesAverage  float64 `xml:"bayesaverage"`
	Rank          []Rank  `xml:"ranks>rank"`
}

type Rank struct {
	Type         string  `xml:"type,attr"`
	ID           int     `xml:"id,attr"`
	Name         string  `xml:"name,attr"`
	Friendly     string  `xml:"friendlyname,attr"`
	Value        int     `xml:"value,attr"`
	BayesAverage float64 `xml:"bayesaverage,attr"`
}

type Poll struct {
	Name       string       `xml:"name,attr"`
	Title      string       `xml:"title,attr"`
	TotalVotes int          `xml:"totalvotes,attr"`
	Results    []PollResult `xml:"results"`
}

type PollResult struct {
	NumPlayers string        `xml:"numplayers,attr,omitempty"`
	Level      string        `xml:"level,attr,omitempty"`
	Values     []ResultValue `xml:"result"`
}

type ResultValue struct {
	Value    string `xml:"value,attr"`
	NumVotes int    `xml:"numvotes,attr"`
}

type PollSummary struct {
	Name    string        `xml:"name,attr"`
	Title   string        `xml:"title,attr"`
	Results []SummaryItem `xml:"result"`
}

type SummaryItem struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
