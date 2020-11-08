package allow

import "encoding/xml"

type LogEntry struct {
	Date          string `csv:"Date"`
	User          string `csv:"User"`
	Path          string `csv:"Path"`
	LocalAddress  string `csv:"Address (Local)"`
	LocalPort     string `csv:"Port (Local)"`
	RemoteAddress string `csv:"Address (Remote)"`
	RemotePort    string `csv:"Port (Remote)"`
	Protocol      string `csv:"Protocol"`
	FilterName    string `csv:"Filter name"`
	FilterID      string `csv:"Filter ID"`
	Direction     string `csv:"Direction"`
	State         string `csv:"State"`
}

type Profile struct {
	XMLName   xml.Name `xml:"root"`
	Timestamp string   `xml:"timestamp,attr"`
	Type      string   `xml:"type,attr"`
	Version   string   `xml:"version,attr"`
	Apps      struct {
		Item []struct {
			Path      string `xml:"path,attr"`
			Timestamp string `xml:"timestamp,attr"`
			IsEnabled string `xml:"is_enabled,attr"`
			IsSilent  string `xml:"is_silent,attr"`
		} `xml:"item"`
	} `xml:"apps"`
	RulesCustom struct {
		Item []RuleItem `xml:"item"`
	} `xml:"rules_custom"`
	RulesConfig string `xml:"rules_config"`
}

type RuleItem struct {
	Name      string `xml:"name,attr"`
	Rule      string `xml:"rule,attr"`
	Protocol  string `xml:"protocol,attr"`
	IsEnabled string `xml:"is_enabled,attr"`
}
