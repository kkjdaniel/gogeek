package thing

type Items struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Type          string        `xml:"type,attr"`
	ID            int           `xml:"id,attr"`
	Name          []Name        `xml:"name"`
	Description   string        `xml:"description"`
	YearPublished IntValue      `xml:"yearpublished"`
	MinPlayers    IntValue      `xml:"minplayers"`
	MaxPlayers    IntValue      `xml:"maxplayers"`
	PlayingTime   IntValue      `xml:"playingtime"`
	MinPlayTime   IntValue      `xml:"minplaytime"`
	MaxPlayTime   IntValue      `xml:"maxplaytime"`
	MinAge        IntValue      `xml:"minage"`
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

type IntValue struct {
	Value int `xml:"value,attr"`
}

type FloatValue struct {
	Value float64 `xml:"value,attr"`
}

type StringValue struct {
	Value string `xml:"value,attr"`
}

type Link struct {
	Type  string `xml:"type,attr"`
	ID    int    `xml:"id,attr"`
	Value string `xml:"value,attr"`
}

type Statistics struct {
	UsersRated    IntValue   `xml:"usersrated"`
	Average       FloatValue `xml:"average"`
	BayesAverage  FloatValue `xml:"bayesaverage"`
	Ranks         []Rank     `xml:"ranks>rank"`
	StdDev        FloatValue `xml:"stddev"`
	Median        IntValue   `xml:"median"`
	Owned         IntValue   `xml:"owned"`
	Trading       IntValue   `xml:"trading"`
	Wanting       IntValue   `xml:"wanting"`
	Wishing       IntValue   `xml:"wishing"`
	NumComments   IntValue   `xml:"numcomments"`
	NumWeights    IntValue   `xml:"numweights"`
	AverageWeight FloatValue `xml:"averageweight"`
}
type Rank struct {
	Type     string `xml:"type,attr"`
	ID       int    `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Friendly string `xml:"friendlyname,attr"`
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
