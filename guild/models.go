package guild

// Guild is a BGG guild (user group) as returned by the guild endpoint.
type Guild struct {
	ID          int      `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Created     string   `xml:"created,attr"`
	Category    string   `xml:"category"`
	Website     string   `xml:"website"`
	Manager     string   `xml:"manager"`
	Description string   `xml:"description"`
	Location    Location `xml:"location"`
	Members     *Members `xml:"members"`
}

// Members is a paginated list of a guild's members.
type Members struct {
	Count   int      `xml:"count,attr"`
	Page    int      `xml:"page,attr"`
	Members []Member `xml:"member"`
}

// Member is a single guild member and the date they joined.
type Member struct {
	Name string `xml:"name,attr"`
	Date string `xml:"date,attr"`
}

// Location is the postal address of a guild.
type Location struct {
	Addr1           string `xml:"addr1"`
	Addr2           string `xml:"addr2"`
	City            string `xml:"city"`
	StateOrProvince string `xml:"stateorprovince"`
	PostalCode      string `xml:"postalcode"`
	Country         string `xml:"country"`
}
