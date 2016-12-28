package persist

import (
	"errors"
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"math/rand"
	"strings"
	"time"
)

type DatastorePersistenceManager struct {
	ctx context.Context
}

const entityTypeName string = "Name"
const entityTypeStats string = "Stats"
const entityTypeRecommendations = "Recommendations"
const filterNameEquals string = "Name = "
const filterGenderEquals string = "Gender ="
const filterRandomNumber string = "Random >="
const filterRecommendationUser string = "Email ="
const filterRecommendationBool string = "Recommended ="

var NoRandomName = errors.New("No random name returned.")
var NoNameFound = errors.New("Requested name not found.")

func NewDatastoreManager(ctx context.Context) *DatastorePersistenceManager {
	return &DatastorePersistenceManager{
		ctx: ctx,
	}
}

func (mgr *DatastorePersistenceManager) GetName(name string, gender names.Gender) (*names.Name, error) {
	query := mgr.buildNameQuery(name, gender)
	result, _, err := mgr.executeGetNameQuery(query)
	if err != nil {
		log.Warningf(mgr.ctx, "action=get_name requested_name=%v gender=%v error=%v", name, gender.GoString(), err)
		return nil, err
	}

	if len(result) == 0 {
		log.Warningf(mgr.ctx, "action=get_name requested_name=%v gender=%v result=no_names_found",
			name, gender.GoString())
		return nil, NoNameFound
	}

	if len(result) > 1 {
		log.Warningf(mgr.ctx, "action=get_name requested_name=%v gender=%v result=multiple_names_found",
			name, gender.GoString())
	}

	return result[0], nil
}

func (mgr *DatastorePersistenceManager) GetRandomName(gender names.Gender) (*names.Name, error) {
	rnd := randomFloat()
	query := datastore.NewQuery(entityTypeName).
		Filter(filterGenderEquals, gender.GoString()).
		Filter(filterRandomNumber, rnd).
		Limit(1)

	results, _, err := mgr.executeGetNameQuery(query)
	if err != nil {
		log.Errorf(mgr.ctx, "Failed to get random name: %v", err)
		return nil, err
	}
	if len(results) == 0 {
		log.Errorf(mgr.ctx, "No random name returned.")
		return nil, NoRandomName
	}
	return results[0], err
}

func (mgr *DatastorePersistenceManager) AddName(name *names.Name) error {
	//Check if there is an existing name in the datastore
	_, keys, err := mgr.executeGetNameQuery(mgr.buildNameQuery(name.Name, name.Gender))
	if err != nil {
		log.Errorf(mgr.ctx, "Action=add_name_to_datastore %v", err)
	}

	dName := newDatastoreName(name)
	var key *datastore.Key

	//If there is an existing key for this name use it, otherwise create a new one.
	if len(keys) > 0 {
		key = keys[0]
	} else {
		key = datastore.NewKey(mgr.ctx, entityTypeName, dName.Key(), 0, nil)
	}

	//Put the name in the datastore.
	key, err = datastore.Put(mgr.ctx, key, dName)
	if err != nil {
		log.Errorf(mgr.ctx, "Failed to put name %v into datastore: %v", name.Name, err)
		return err
	}

	//Add the stats
	err = mgr.addStatsToDatastore(mgr.ctx, key, name.Stats)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *DatastorePersistenceManager) GetRecommendedNames(usr, partner *user.User) ([]*names.Name, error) {

	//Get all of the recommendations for this user
	uRec, uKeys, err := mgr.getUserRecommendations(mgr.getRecommendationQuery(usr, false))
	if err != nil {
		log.Errorf(mgr.ctx, "action=GetRecommendedNames error=%v", err)
		return nil, err
	}

	//Add them to a map
	nameMap := make(map[string]*decision.Recommendation)
	for index, key := range uKeys {
		nameMap[key.Parent().StringID()] = uRec[index]
	}

	//Get all of the recommendations for the user
	_, pKeys, err := mgr.getUserRecommendations(mgr.getRecommendationQuery(partner, true))
	if err != nil {
		log.Errorf(mgr.ctx, "action=GetRecommendedNames error=%v", err)
		return nil, err
	}

	delta := []string{}
	//determine if this is a recommendation that the partner has made but the user hasn't decided on.
	for _, key := range pKeys {
		if nameMap[key.Parent().StringID()] == nil {
			delta = append(delta, key.Parent().StringID())
		}
	}

	result := []*names.Name{}
	for _, stringKey := range delta {
		parts := strings.Split(stringKey, "::")
		nameVals, _, err := mgr.executeGetNameQuery(mgr.buildNameQuery(parts[0], names.GetGender(parts[1])))
		if err != nil {
			log.Errorf(mgr.ctx, "action=GetRecommendedNames error=%v", err)
		}
		for _, name := range nameVals {
			result = append(result, name)
		}
	}

	return result, nil
}

func (mgr *DatastorePersistenceManager) getRecommendationQuery(usr *user.User, approvedOnly bool) *datastore.Query {
	q := datastore.NewQuery(entityTypeRecommendations).
		Filter(filterRecommendationUser, usr.Email)
	if approvedOnly {
		q = q.Filter(filterRecommendationBool, true)
	}
	return q
}

func (mgr *DatastorePersistenceManager) getUserRecommendations(
	query *datastore.Query) ([]*decision.Recommendation, []*datastore.Key, error) {
	recommendations := []*decision.Recommendation{}
	keys, err := query.GetAll(mgr.ctx, &recommendations)
	return recommendations, keys, err
}

func (mgr *DatastorePersistenceManager) UpdateDecision(name *names.Name, d *decision.Recommendation) error {
	_, keys, err := mgr.executeGetNameQuery(mgr.buildNameQuery(name.Name, name.Gender))

	if err != nil {
		log.Errorf(mgr.ctx, "action=UpdateDecision error=%v", err.Error())
		return err
	}

	if len(keys) > 0 {
		for _, key := range keys {
			decisionKeys, err := mgr.getDecisions(key, d)
			if err != nil {
				return err
			}

			if len(decisionKeys) > 0 {
				for _, decisionKey := range decisionKeys {
					_, err := datastore.Put(mgr.ctx, decisionKey, d)
					if err != nil {
						log.Errorf(mgr.ctx, "action=UpdateDecision error=%v", err)
						return err
					}
				}
			} else {
				err := mgr.writeNewDecision(key, d)
				if err != nil {
					log.Errorf(mgr.ctx, "action=UpdateDecision error=%v", err)
					return err
				}
			}
		}

	}
	return nil
}

func (mgr *DatastorePersistenceManager) getDecisions(ancestor *datastore.Key, decision *decision.Recommendation) ([]*datastore.Key, error) {
	//get decisions
	query := datastore.NewQuery(entityTypeRecommendations).
		Ancestor(ancestor).Filter(filterRecommendationUser, decision.Email).KeysOnly()
	results := make([]interface{}, 0)
	decisionKeys, err := query.GetAll(mgr.ctx, results)
	if err != nil {
		log.Errorf(mgr.ctx, "action=UpdateDecision error=%v", err.Error())
		return nil, err
	}
	return decisionKeys, nil
}

func (mgr *DatastorePersistenceManager) writeNewDecision(parent *datastore.Key, decision *decision.Recommendation) error {
	decisionKey := datastore.NewIncompleteKey(mgr.ctx, entityTypeRecommendations, parent)
	_, err := datastore.Put(mgr.ctx, decisionKey, decision)
	if err != nil {
		log.Errorf(mgr.ctx, "action=writeNewDecision error=%v", err)
		return err
	}
	return nil
}

func (mgr *DatastorePersistenceManager) buildNameQuery(name string, gender names.Gender) *datastore.Query {
	return datastore.NewQuery(entityTypeName).
		Filter(filterNameEquals, name).
		Filter(filterGenderEquals, gender.GoString())
}

func (mgr *DatastorePersistenceManager) executeGetNameQuery(query *datastore.Query) ([]*names.Name, []*datastore.Key, error) {

	results := []*datastoreName{}
	keys, err := query.GetAll(mgr.ctx, &results)

	if err != nil {
		return nil, nil, errors.New("Failed to get name: " + err.Error())
	}

	result := []*names.Name{}
	for index, name := range results {
		out := &name.Name
		out.Gender = names.GetGender(name.Gender)
		mgr.getStatsForKey(out, keys[index])
		result = append(result, out)
	}

	return result, keys, nil
}

func (mgr *DatastorePersistenceManager) getStatsForKey(name *names.Name, key *datastore.Key) (*names.Name, error) {

	statsQuery := datastore.NewQuery(entityTypeStats).Ancestor(key)
	stats := []*names.Stat{}

	_, err := statsQuery.GetAll(mgr.ctx, &stats)
	if err != nil {
		return nil, errors.New("Failed to retrieve stats: " + err.Error())
	}

	for _, stat := range stats {
		name.AddStat(stat)
	}
	return name, nil
}

func (mgr *DatastorePersistenceManager) addStatsToDatastore(ctx context.Context,
	parent *datastore.Key,
	stats map[int]*names.Stat) error {

	for _, stat := range stats {
		key := datastore.NewIncompleteKey(ctx, entityTypeStats, parent)
		_, err := datastore.Put(ctx, key, stat)
		if err != nil {
			return err
		}
	}
	return nil
}

func randomFloat() float32 {
	return rand.New(rand.NewSource(time.Now().Unix())).Float32()
}
