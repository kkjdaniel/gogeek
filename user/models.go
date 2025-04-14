package user

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
}

type ValueField struct {
	Value string `xml:"value,attr"`
}

type IntValueField struct {
	Value int `xml:"value,attr"`
}

type Buddies struct {
	Total int     `xml:"total,attr"`
	Page  int     `xml:"page,attr"`
	Buddy []Buddy `xml:"buddy"`
}

type Buddy struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type Guilds struct {
	Total int     `xml:"total,attr"`
	Page  int     `xml:"page,attr"`
	Guild []Guild `xml:"guild"`
}

type Guild struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type Top struct {
	Domain string    `xml:"domain,attr"`
	Items  []TopItem `xml:"item"`
}

type TopItem struct {
	Rank int    `xml:"rank,attr"`
	Type string `xml:"type,attr"`
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}
