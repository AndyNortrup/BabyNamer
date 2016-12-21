package babynamer

import (
	"net/http"
	"net/url"

	"github.com/AndyNortrup/baby-namer/names"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func namesPage(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	ctx := appengine.NewContext(r)
	username := user.Current(ctx)
	gen := NewNameGenerator(ctx, username.String())

	name, decision := getQueryParam(r.URL)

	recordDecision(name, decision, username, ctx)

	newName, err := gen.getName(name)
	if err != nil {
		log.Errorf(ctx, "Error getting name: %v", err)
	}

	renderNamePageTemplate(newName, ctx, w)

	gen.getRandomName(names.FemaleFilter)
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

func renderNamePageTemplate(newName *names.Name, ctx context.Context, w http.ResponseWriter) {
	t, err := getNameTemplate()
	if err != nil {
		log.Errorf(ctx, "Failed to parse template file: %v", err)
	}

	err = t.Execute(w, newName)
	if err != nil {
		log.Errorf(ctx, "Failed to Execute template: %v\tName: %v", err, newName)
	}
}
