package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func handleShortList(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	list := NewShortList(ctx)
	t, err := getShortListTemplate()

	if err != nil {
		log.Errorf(ctx, "Unable to get short list template: %v", err)
	}

	err = t.Execute(w, list)
	if err != nil {
		log.Errorf(ctx, "Failed to execute short list template: %v", err)
	}
}
