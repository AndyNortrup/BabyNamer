package main

import (
	"golang.org/x/net/context"

	"errors"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type DatastoreNameManager struct {
	name, username string
	ctx            context.Context
}

func (u *DatastoreNameManager) getNameFromDatastore() (*NameDetails, *datastore.Key, error) {
	//Get the name from the datastore
	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("Name =", u.name)

	for t := query.Run(u.ctx); ; {
		details := &NameDetails{}
		key, err := t.Next(details)

		if err == datastore.Done {
			log.Infof(u.ctx, "Couldn't find name in datastore.")
			return nil, nil, err
		}
		if err != nil {
			log.Warningf(u.ctx, "Error retriving Name Details %#v", err)
			return nil, nil, err
		} else {
			return details, key, nil
		}
	}

	return nil, nil, errors.New("Unable to locate name in Datastore.")
}

func (u *DatastoreNameManager) updateNameRecommendations(
	decision bool) error {

	record, key, err := u.getNameFromDatastore()

	if err != nil {
		return errors.New("Unable to locate name in datastore.")
	}

	err = u.recordDecision(key, decision, record)
	if err != nil {
		log.Warningf(u.ctx, "Error updating name with decision.")
		return err
	}
	log.Infof(u.ctx, "Name record updated.")
	return nil
}

func (u *DatastoreNameManager) recordDecision(key *datastore.Key,
	decision bool,
	details *NameDetails) error {

	err := checkDuplicateDecision(details, u.username)

	if err != nil {
		return err
	}

	if decision {
		//Check if this has already been updated
		if details.RecommendedBy == "" {
			details.RecommendedBy = u.username
		} else {
			details.ApprovedBy = u.username
		}
	} else {
		details.RejectedBy = u.username
	}

	_, err = datastore.Put(u.ctx, key, details)
	return err
}

//checkDuplicateDecision prevents the same user for redering judgement on a name more than once.
func checkDuplicateDecision(details *NameDetails, username string) error {
	//Make sure that this user hasn't already recommended / approved
	// rejected this name.
	if details.RecommendedBy == username {
		return NewDuplicateDecisionError(username, "recommended")
	} else if details.ApprovedBy == username {
		return NewDuplicateDecisionError(username, "approved")
	} else if details.RejectedBy == username {
		return NewDuplicateDecisionError(username, "rejected")
	}
	return nil
}

type DuplicateDecisionError struct {
	message string
}

func NewDuplicateDecisionError(user, action string) *DuplicateDecisionError {
	return &DuplicateDecisionError{
		message: "User " + user + " already " + action + "this name",
	}
}

func (err DuplicateDecisionError) Error() string {
	return err.message
}
