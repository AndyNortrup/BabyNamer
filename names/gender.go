package names

type Gender struct {
	female bool
}

func GetGender(str string) Gender {
	if str == "F" {
		return Female
	} else {
		return Male
	}
}

func (g *Gender) GoString() string {
	if g.female {
		return "F"
	} else {
		return "M"
	}
}

var Male Gender = Gender{female: false}
var Female Gender = Gender{female: true}
