package babynamer

import (
	"golang.org/x/net/context"

	"errors"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type DatastoreNameManager struct {
	username string
	ctx      context.Context
}

var DuplicateNameError = errors.New("Name already exists in datastore.")

func NewDatastoreNameManager(ctx context.Context, username string) *DatastoreNameManager {
	return &DatastoreNameManager{
		username: username,
		ctx:      ctx,
	}
}

func (u *DatastoreNameManager) getNameFromDatastore(name string) (*NameDetails, *datastore.Key, error) {
	//Get the name from the datastore
	log.Debugf(u.ctx, "Getting name from datastore: %v", name)
	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("Name =", name)

	for t := query.Run(u.ctx); ; {
		details := &NameDetails{}
		key, err := t.Next(details)

		if err == datastore.Done {
			log.Infof(u.ctx, "Couldn't find name in datastore: %v", name)
			return nil, nil, err
		}
		if err != nil {
			log.Warningf(u.ctx, "Error retriving Name Details %#v", err)
			return nil, nil, err
		} else {
			return details, key, nil
		}
	}
}

func (u *DatastoreNameManager) updateNameRecommendations(name string,
	decision bool) error {

	record, key, err := u.getNameFromDatastore(name)

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
		details.Approve(u.username)
	} else {
		details.Reject(u.username)
	}

	_, err = datastore.Put(u.ctx, key, details)
	return err
}

func (u *DatastoreNameManager) addNameToStore(details *NameDetails) error {

	if details.Name == "" {
		return errors.New("Can't write blank name to datastore.")
	}

	//Check if the name already exists by searching for it.
	key, _ := u.getKeyForName(details)
	if key != nil {
		return DuplicateNameError
	}

	key = datastore.NewIncompleteKey(u.ctx, EntityTypeNameDetails, nil)
	if _, err := datastore.Put(u.ctx, key, details); err != nil {
		log.Warningf(u.ctx, "Error writing name to datastore: %v", err)
		return err
	} else {
		log.Infof(u.ctx, "Added name to the datastore: %v", details.Name)
	}
	return nil
}

//getKeyForName returns the datastore key for a given name value. Use this to
// get a key so you update a name's information rather than write a new copy.
func (u *DatastoreNameManager) getKeyForName(details *NameDetails) (*datastore.Key, error) {
	t := datastore.NewQuery(EntityTypeNameDetails).
		Filter("Name =", details.Name).
		Run(u.ctx)

	for {
		results := &NameDetails{}
		key, err := t.Next(results)

		if err == datastore.Done {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}
		return key, nil
	}
}
