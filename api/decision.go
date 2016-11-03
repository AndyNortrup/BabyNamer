package main

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func namesPage(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	ctx := appengine.NewContext(r)
	user := user.Current(ctx)

	gen := NewNameGenerator(ctx)

	name := r.URL.Query().Get("name")
	var decision bool

	if r.URL.Query().Get("decision") == "yes" {
		decision = true
	} else {
		decision = false
	}

	updateNameRecommendations(name, user.Email, decision, ctx)

	newName, err := gen.getName(name)
	if err != nil {
		log.Errorf(ctx, "Error getting name: %v", err)
	}

	t, err := template.ParseFiles("templates/name-suggestor.html")
	if err != nil {
		log.Errorf(ctx, "Failed to parse template file: %v", err)
	}

	err = t.Execute(w, newName)
	if err != nil {
		log.Errorf(ctx, "Failed to Execute template: %v", err)
	}
}
