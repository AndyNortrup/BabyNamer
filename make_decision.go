package babynamer

import (
	"net/http"
	"net/url"

	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/persistance"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"github.com/AndyNortrup/baby-namer/settings"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func namesPage(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	ctx := appengine.NewContext(r)
	username := user.Current(ctx)

	name, rec := getDecisionFromURL(r.URL, username)
	recordDecision(name, rec, ctx)

	partner, _ := settings.GetPartner(ctx, username)
	if partner == nil {
		sp := NewSettingsPage(ctx)
		sp.Render(w)
	} else {
		//Create suggestion page and render it.
		sp := NewSuggestionPage(username, persist.NewDatastoreManager(ctx), ctx, name)
		sp.getName()
		sp.render(w)
	}
}

func recordDecision(name *names.Name, rec *decision.Recommendation, ctx context.Context) {
	nameMgr := persist.NewDatastoreManager(ctx)
	err := nameMgr.UpdateDecision(name, rec)
	if err != nil {
		log.Errorf(ctx, "action=recordDecision error=%v", err)
	} else {
		log.Infof(ctx, "action=recordDecision name=%v recommendation=%v", name, rec)
	}

}

func getDecisionFromURL(url *url.URL, usr *user.User) (*names.Name, *decision.Recommendation) {

	name := names.NewName(
		url.Query().Get("name"),
		names.GetGender(url.Query().Get("gender")))

	var d bool
	if url.Query().Get("decision") == "yes" {
		d = true
	} else {
		d = false
	}
	dec := decision.NewRecommendation(usr, d)

	return name, dec
}
