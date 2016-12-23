package babynamer

import (
	"net/http"
	"net/url"

	"github.com/AndyNortrup/baby-namer/persistance"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func namesPage(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	ctx := appengine.NewContext(r)
	username := user.Current(ctx)

	name, decision := getQueryParam(r.URL)
	recordDecision(name, decision, username, ctx)

	//Create suggestion page and render it.
	sp := NewSuggestionPage(username, persist.NewDatastoreManager(ctx), ctx)
	sp.getName()
	sp.render(w)
}

func recordDecision(name string, decision bool, username *user.User, ctx context.Context) {
	nameMgr := NewDatastoreNameManager(ctx, username.String())

	if name != "" {
		nameMgr.updateNameRecommendations(name, decision)
	}

}

func getQueryParam(url *url.URL) (string, bool) {

	name := url.Query().Get("name")
	var decision bool

	if url.Query().Get("decision") == "yes" {
		decision = true
	} else {
		decision = false
	}
	return name, decision
}
