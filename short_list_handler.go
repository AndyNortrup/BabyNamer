package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type ShortListPage struct {
	NameList ShortList
}

func handleShortList(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	list, err := getShortList(ctx)
	if err != nil {
		log.Errorf(ctx, "Failed to list get short list: %v", err)
	}

	page := ShortListPage{
		NameList: list,
	}
	t, err := getShortListTemplate()

	if err != nil {
		log.Errorf(ctx, "Unable to get short list template: %v", err)
	}

	err = t.Execute(w, page)
	if err != nil {
		log.Errorf(ctx, "Failed to execute short list template: %v", err)
	}
}
