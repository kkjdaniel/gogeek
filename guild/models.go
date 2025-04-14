package guild

type Guild struct {
	ID          int      `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Created     string   `xml:"created,attr"`
	Category    string   `xml:"category"`
	Website     string   `xml:"website"`
	Manager     string   `xml:"manager"`
	Description string   `xml:"description"`
	Location    Location `xml:"location"`
}

type Location struct {
	Addr1           string `xml:"addr1"`
	Addr2           string `xml:"addr2"`
	City            string `xml:"city"`
	StateOrProvince string `xml:"stateorprovince"`
	PostalCode      string `xml:"postalcode"`
	Country         string `xml:"country"`
}
