package main

type specialReport struct {
	UUID                   string                 `json:"uuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	PrefLabel              string                 `json:"prefLabel"`
	Type                   string                 `json:"type"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	Uuids []string `json:"uuids,omitempty"`
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
