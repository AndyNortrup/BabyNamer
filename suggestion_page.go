package babynamer

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/persistance"
	"github.com/AndyNortrup/baby-namer/settings"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"net/http"
)

type SuggestionPage struct {
	ctx             context.Context
	Name            *names.Name
	MostOccurrences *names.Stat
	HighestRank     *names.Stat
	Gender          string
	usr             *user.User
	data            persist.DataManager
	lastName        *names.Name
	IsRecommended   bool
}

func NewSuggestionPage(u *user.User,
	data persist.DataManager,
	ctx context.Context,
	lastName *names.Name) *SuggestionPage {
	sp := &SuggestionPage{usr: u, data: data, ctx: ctx, lastName: lastName}
	return sp
}

func (sp *SuggestionPage) getName() {

	name := sp.recommendedName()

	if name == nil {
		//Todo: better management of errors to the front end.
		name, _ = sp.randomName()
	}
	sp.addDetailsToPage(name)

}

func (sp *SuggestionPage) recommendedName() *names.Name {
	sp.IsRecommended = false

	self := user.Current(sp.ctx)
	partner, err := settings.GetPartner(sp.ctx, self)
	if err != nil {
		return nil
	}

	uPartner := &user.User{Email: partner.PartnerEmail}

	recNames, err := sp.data.GetPartnerRecommendedNames(self, uPartner)
	if err != nil || len(recNames) == 0 {
		return nil
	}

	if sp.lastName != nil {
		for _, n := range recNames {
			if n.Name != sp.lastName.Name {
				sp.IsRecommended = true
				return n
			}
		}
	} else {
		sp.IsRecommended = true
		return recNames[0]
	}

	return nil
}

func (sp *SuggestionPage) randomName() (*names.Name, error) {
	needName := true

	for needName {
		name, err := sp.data.GetRandomName(names.Female)
		if err != nil {
			log.Errorf(sp.ctx, err.Error())
			return nil, err
		}
		recs, err := sp.data.GetNameRecommendations(sp.usr, name)
		if err != nil {
			return nil, err
		}
		if len(recs) == 0 {
			return name, err
		}
	}
	return nil, nil
}

func (sp *SuggestionPage) addDetailsToPage(name *names.Name) {
	sp.Name = name
	sp.HighestRank = name.MostOccurrences()
	sp.MostOccurrences = name.MostOccurrences()
	sp.Gender = name.Gender.GoString()
}

func (sp *SuggestionPage) render(w http.ResponseWriter) {

	t, err := getNameTemplate()
	if err != nil {
		log.Errorf(sp.ctx, "Failed to parse template file: %v", err)
	}

	err = t.Execute(w, sp)
	if err != nil {
		log.Errorf(sp.ctx, "Failed to Execute template: %v\tName: %v", err, sp)
	}

}
