package shodan

type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	Https        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

type Location struct {
	City         string  `json:"city"`
	Lat          float32 `json:"lat"`
	Long         float32 `json:"long"`
	RegionCode   string  `json:"region_code"`
	AreaCode     string  `json:"area_code"`
	CountryCode3 string  `json:"country_code3"`
	CountryName  string  `json:"country_name"`
	PostalCode   string  `json:"postal_code"`
	DMACode      string  `json:"dma_code"`
	CountryCode  string  `json:"country_code"`
}

type Host struct {
	OS           string   `json:"os"`
	Timestamp    string   `json:"timestamp"`
	ISP          string   `json:"isp"`
	ASN          string   `json:"asn"`
	Hostnames    []string `json:"hostnames"`
	Location     Location `json:"location"`
	IP           string   `json:"ip"`
	Domains      []string `json:"domains"`
	Organization string   `json:"organization"`
	Data         string   `json:"data"`
	Port         string   `json:"port"`
	IPString     string   `json:"ip_string"`
}

type SearchResult struct {
	Matches []Host `json:"matches"`
}
