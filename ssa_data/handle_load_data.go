package ssa_data

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
)

func HandleLoadData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)

	names := loadNames()
	err := addNamesToDatastore(ctx, names)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func addNamesToDatastore(ctx context.Context, names map[string]*Name) error {
	for _, value := range names {
		key := datastore.NewKey(ctx, entityTypeName, value.makeNameKey(), 0, nil)
		key, err := datastore.Put(ctx, key, value)
		if err != nil {
			log.Errorf(ctx, "Failed to put name %v into datastore: %v", value.Name, err)
			return err
		}

		err = addStatsToDatastore(ctx, key, value.Stats)
		if err != nil {
			return err
		}
	}
	return nil
}

func addStatsToDatastore(ctx context.Context, parent *datastore.Key, stats map[int]*Stat) error {

	for _, stat := range stats {
		key := datastore.NewIncompleteKey(ctx, entityTypeStats, parent)
		_, err := datastore.Put(ctx, key, stat)
		if err != nil {
			return err
		}
	}
	return nil
}
