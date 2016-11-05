package main

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func updateNameRecommendations(name, username string,
	decision bool,
	ctx context.Context) error {

	//Get the name from the datastore
	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("Name =", name)

	for t := query.Run(ctx); ; {
		details := &NameDetails{}
		key, err := t.Next(details)

		if err == datastore.Done {
			log.Infof(ctx, "Couldn't find name in datastore.")
			break
		}
		if err != nil {
			log.Warningf(ctx, "Error retriving Name Details %#v", err)
			return err
		} else {
			err := recordDecision(key, decision, details, username, ctx)
			if err != nil {
				log.Warningf(ctx, "Error updating name with decision.")
				return err
			}
			log.Infof(ctx, "Name record updated.")
			break
		}
	}
	return nil
}

func recordDecision(key *datastore.Key,
	decision bool,
	details *NameDetails,
	username string,
	ctx context.Context) error {

	if decision {
		//Check if this has already been updated
		if details.RecommendedBy == "" {
			details.RecommendedBy = username
		} else {
			details.ApprovedBy = username
		}
	} else {
		details.RejectedBy = username
	}

	_, err := datastore.Put(ctx, key, details)
	return err
}
