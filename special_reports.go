package main

type specialReport struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier"`
	Type          string `json:"type"`
}

type specialReportVariation struct {
	Name      string   `json:"name"`
	Weight    string   `json:"weight"`
	Case      string   `json:"case"`
	Accent    string   `json:"accent"`
	Languages []string `json:"languages"`
}

type specialReportLink struct {
	APIURL string `json:"apiUrl"`
}
