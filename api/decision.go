package main

import (
	"net/http"
	"net/url"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func namesPage(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	ctx := appengine.NewContext(r)
	username := user.Current(ctx)

	gen := NewNameGenerator(ctx)
	name, decision := getQueryParam(r.URL)

	nameMgr := DatastoreNameManager{
		username: username.String(),
		name:     name,
	}
	nameMgr.updateNameRecommendations(decision)

	newName, err := gen.getName(name)
	if err != nil {
		log.Errorf(ctx, "Error getting name: %v", err)
	}

	t, err := getNameTemplate()
	if err != nil {
		log.Errorf(ctx, "Failed to parse template file: %v", err)
	}

	err = t.Execute(w, newName)
	if err != nil {
		log.Errorf(ctx, "Failed to Execute template: %v", err)
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
