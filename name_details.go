package babynamer

import "github.com/AndyNortrup/baby-namer/usage"

type NameDetails struct {
	Name    string        `xml:"name_detail>name"`
	Gender  string        `xml:"name_detail>gender"`
	Usages  []usage.Usage `xml:"name_detail>usage" datastore:-`
	Origins []*usage.NameOrigin
	Decision
}
