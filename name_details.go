package main

type NameDetails struct {
	Name          string  `xml:"name_detail>name"`
	Gender        string  `xml:"name_detail>gender"`
	Usages        []Usage `xml:"name_detail>usages"`
	RecommendedBy string
	ApprovedBy    string
	RejectedBy    string
}
