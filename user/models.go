package user

// User is a BGG user profile returned by the user endpoint.
type User struct {
	ID               int           `xml:"id,attr"`
	Name             string        `xml:"name,attr"`
	FirstName        ValueField    `xml:"firstname"`
	LastName         ValueField    `xml:"lastname"`
	AvatarLink       ValueField    `xml:"avatarlink"`
	YearRegistered   IntValueField `xml:"yearregistered"`
	LastLogin        ValueField    `xml:"lastlogin"`
	StateOrProvince  ValueField    `xml:"stateorprovince"`
	Country          ValueField    `xml:"country"`
	WebAddress       ValueField    `xml:"webaddress"`
	XboxAccount      ValueField    `xml:"xboxaccount"`
	WiiAccount       ValueField    `xml:"wiiaccount"`
	PSNAccount       ValueField    `xml:"psnaccount"`
	BattleNetAccount ValueField    `xml:"battlenetaccount"`
	SteamAccount     ValueField    `xml:"steamaccount"`
	TradeRating      IntValueField `xml:"traderating"`
	Buddies          Buddies       `xml:"buddies"`
	Guilds           Guilds        `xml:"guilds"`
	Top              Top           `xml:"top"`
	Hot              Hot           `xml:"hot"`
}

// ValueField holds a string XML value attribute.
type ValueField struct {
	Value string `xml:"value,attr"`
}

// IntValueField holds an integer XML value attribute.
type IntValueField struct {
	Value int `xml:"value,attr"`
}

// Buddies is the paginated list of a user's buddies.
type Buddies struct {
	Total int     `xml:"total,attr"`
	Page  int     `xml:"page,attr"`
	Buddy []Buddy `xml:"buddy"`
}

// Buddy is a single user in a user's buddy list.
type Buddy struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// Guilds is the paginated list of guilds a user belongs to.
type Guilds struct {
	Total int     `xml:"total,attr"`
	Page  int     `xml:"page,attr"`
	Guild []Guild `xml:"guild"`
}

// Guild is a single guild a user belongs to.
type Guild struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// Top is a user's ranked "top" list for a domain (e.g. boardgame).
type Top struct {
	Domain string    `xml:"domain,attr"`
	Items  []TopItem `xml:"item"`
}

// Hot is a user's ranked "hot" list for a domain (e.g. boardgame).
type Hot struct {
	Domain string    `xml:"domain,attr"`
	Items  []TopItem `xml:"item"`
}

// TopItem is a single ranked entry in a user's top or hot list.
type TopItem struct {
	Rank int    `xml:"rank,attr"`
	Type string `xml:"type,attr"`
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}
