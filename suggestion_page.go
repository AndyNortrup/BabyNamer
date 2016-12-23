package babynamer

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/persistance"
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
	u               *user.User
	data            persist.DataManager
}

func NewSuggestionPage(u *user.User, data persist.DataManager, ctx context.Context) *SuggestionPage {
	sp := &SuggestionPage{u: u, data: data, ctx: ctx}
	return sp
}

func (sp *SuggestionPage) getName() {
	name := sp.randomName()
	sp.addDetailsToPage(name)

}

func (sp *SuggestionPage) randomName() *names.Name {
	name, err := sp.data.GetRandomName(names.Female)
	if err != nil {
		log.Errorf(sp.ctx, err.Error())
		return nil
	}
	return name
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
