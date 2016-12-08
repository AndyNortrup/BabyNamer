package usage

type Usage struct {
	UsageFull   string `xml:"usage>usage_full"`
	UsageGender string `xml:"usage>usage_gender"`
}
