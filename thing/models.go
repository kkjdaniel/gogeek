package thing

// Items is the root element returned by the thing endpoint, wrapping the list of requested items.
type Items struct {
	Items []Item `xml:"item"`
}

// Item is a single thing (board game, expansion, accessory, etc.) returned by the thing endpoint.
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

// Name is a primary or alternate name of an item.
type Name struct {
	Type      string `xml:"type,attr"`
	SortIndex int    `xml:"sortindex,attr"`
	Value     string `xml:"value,attr"`
}

// IntValue holds an integer XML value attribute.
type IntValue struct {
	Value int `xml:"value,attr"`
}

// FloatValue holds a floating-point XML value attribute.
type FloatValue struct {
	Value float64 `xml:"value,attr"`
}

// StringValue holds a string XML value attribute.
type StringValue struct {
	Value string `xml:"value,attr"`
}

// Link relates an item to another BGG entity, such as a category, mechanic, designer, or publisher.
type Link struct {
	Type  string `xml:"type,attr"`
	ID    int    `xml:"id,attr"`
	Value string `xml:"value,attr"`
}

// Statistics holds the community rating statistics for an item, such as average rating, rankings, and ownership counts.
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

// Rank is an item's position in a BGG ranking list, such as the overall board game rank or a family rank.
type Rank struct {
	Type         string `xml:"type,attr"`
	ID           int    `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	Friendly     string `xml:"friendlyname,attr"`
	Value        string `xml:"value,attr"`
	BayesAverage string `xml:"bayesaverage,attr"`
}

// Poll is a community poll for an item, such as suggested player counts, player age, or language dependence.
type Poll struct {
	Name       string       `xml:"name,attr"`
	Title      string       `xml:"title,attr"`
	TotalVotes int          `xml:"totalvotes,attr"`
	Results    []PollResult `xml:"results"`
}

// PollResult is a group of vote tallies for one poll option, such as a specific player count.
type PollResult struct {
	NumPlayers string        `xml:"numplayers,attr,omitempty"`
	Level      string        `xml:"level,attr,omitempty"`
	Values     []ResultValue `xml:"result"`
}

// ResultValue is a single poll answer and the number of votes it received.
type ResultValue struct {
	Value    string `xml:"value,attr"`
	NumVotes int    `xml:"numvotes,attr"`
}

// PollSummary is BGG's summarised result of a community poll, such as the best player count.
type PollSummary struct {
	Name    string        `xml:"name,attr"`
	Title   string        `xml:"title,attr"`
	Results []SummaryItem `xml:"result"`
}

// SummaryItem is a single named result within a poll summary, such as "bestwith".
type SummaryItem struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
