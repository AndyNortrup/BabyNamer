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
	"time"
)

//TODO: Replace datastore with nds
type DatastorePersistenceManager struct {
	ctx context.Context
}

const entityTypeName string = "Name"
const entityTypeStats string = "Stats"
const entityTypeRecommendations string = "Recommendations"

const filterNameEquals string = "Name = "
const filterGenderEquals string = "Gender ="
const filterRandomNumber string = "Random >="
const filterRecommendationUser string = "Email ="
const filterRecommendationBool string = "Recommended ="
const filterStatsNameKey string = "NameKey ="
const filterStatYear string = "Year ="

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

//AddName writes a name to the datastore. It always overwrites a previous entity with the same key (Name::Gender)
func (mgr *DatastorePersistenceManager) AddName(name *names.Name) error {
	dName := newDatastoreName(name)
	key := datastore.NewKey(mgr.ctx, entityTypeName, dName.Key(), 0, nil)

	//Put the name in the datastore.
	_, err := datastore.Put(mgr.ctx, key, dName)
	if err != nil {
		log.Errorf(mgr.ctx, "Failed to put name %v into datastore: %v", name.Name, err)
		return err
	}

	//Add the stats
	for _, stat := range name.Stats {
		err = mgr.addStatToDatastore(name, stat)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mgr *DatastorePersistenceManager) GetPartnerRecommendedNames(usr, partner *user.User) ([]*names.Name, error) {

	//Get all of the recommendations for this user
	uRec, err := mgr.getUserRecommendations(mgr.getRecommendationQuery(usr, false))
	if err != nil {
		log.Errorf(mgr.ctx, "action=GetPartnerRecommendedNames error=%v", err)
		return nil, err
	}

	//Add them to a map
	nameMap := make(map[string]*decision.Recommendation)
	for _, key := range uRec {
		name := names.NewName(key.Name, names.GetGender(key.Gender))
		nameMap[name.Key()] = key
	}

	//Get all of the recommendations for the user
	uRec, err = mgr.getUserRecommendations(mgr.getRecommendationQuery(partner, true))
	if err != nil {
		log.Errorf(mgr.ctx, "action=GetPartnerRecommendedNames error=%v", err)
		return nil, err
	}

	delta := []*names.Name{}
	//determine if this is a recommendation that the partner has made but the user hasn't decided on.
	for _, key := range uRec {
		name := names.NewName(key.Name, names.GetGender(key.Gender))
		if nameMap[name.Key()] == nil {
			delta = append(delta, name)
		}
	}

	result := make([]*names.Name, len(delta))
	for index, name := range delta {
		nameVals, _, err := mgr.executeGetNameQuery(mgr.buildNameQuery(name.Name, name.Gender))
		if err != nil {
			log.Errorf(mgr.ctx, "action=GetPartnerRecommendedNames error=%v", err)
		}
		for _, fullName := range nameVals {
			result[index] = fullName
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
	query *datastore.Query) ([]*decision.Recommendation, error) {
	recommendations := []*decision.Recommendation{}
	_, err := query.GetAll(mgr.ctx, &recommendations)
	return recommendations, err
}

//UpdateDecision records the decision in the entity type Recommendations with a string key Name::Gender::UserEmail
// Mary::F::andrew.nortrup@gmail.com.  This overwrites any previous decisions by this user for the name.
func (mgr *DatastorePersistenceManager) UpdateDecision(rec *decision.Recommendation) error {
	key := datastore.NewKey(mgr.ctx,
		entityTypeRecommendations,
		rec.Name+"::"+rec.Gender+"::"+rec.Email,
		0,
		nil)

	_, err := datastore.Put(mgr.ctx, key, rec)
	if err != nil {
		log.Errorf(mgr.ctx, "action=UpdateDeecision error=%v", err)
		return err
	}

	return nil
}

func (mgr *DatastorePersistenceManager) getDecisions(rec *decision.Recommendation) ([]*datastore.Key, error) {
	//get decisions
	query := datastore.NewQuery(entityTypeRecommendations).
		Filter(filterRecommendationUser, rec.Email).
		Filter(filterNameEquals, rec.Name).Filter(filterGenderEquals, rec.Gender).
		KeysOnly()

	results := make([]interface{}, 0)
	decisionKeys, err := query.GetAll(mgr.ctx, results)
	if err != nil {
		log.Errorf(mgr.ctx, "action=UpdateDecision error=%v", err.Error())
		return nil, err
	}
	return decisionKeys, nil
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
		_, err := mgr.getStatsForKey(out, keys[index])
		if err != nil {
			return nil, nil, err
		}
		result = append(result, out)
	}

	return result, keys, nil
}

func (mgr *DatastorePersistenceManager) getStatsForKey(name *names.Name, key *datastore.Key) (*names.Name, error) {

	statsQuery := datastore.NewQuery(entityTypeStats).Filter(filterStatsNameKey, name.Key())
	stats := []*datastoreStat{}

	_, err := statsQuery.GetAll(mgr.ctx, &stats)
	if err != nil {
		log.Errorf(mgr.ctx, "action=getStatsForKey error=err")
		return nil, errors.New("Failed to retrieve stats: " + err.Error())
	}

	for _, stat := range stats {
		name.AddStat(&stat.Stat)
	}
	return name, nil
}

//AddStatsToDataStore adds a instance of a stat to the datastore.  The parent names.Name object is to generate a
// unique key that is Name::Gender::Year (Mary::F::1889)
func (mgr *DatastorePersistenceManager) addStatToDatastore(name *names.Name, stat *names.Stat) error {
	dStat := NewDatastoreStat(name, *stat)
	return mgr.putStat(dStat.newStatKey(mgr.ctx), dStat)
}

func (mgr *DatastorePersistenceManager) putStat(key *datastore.Key, dStat *datastoreStat) error {
	_, err := datastore.Put(mgr.ctx, key, dStat)
	if err != nil {
		log.Errorf(mgr.ctx, "action=addStatToDatastore error=%v", err)
		return err
	}
	return nil
}

func randomFloat() float32 {
	return rand.New(rand.NewSource(time.Now().Unix())).Float32()
}

func (mgr *DatastorePersistenceManager) GetNameRecommendations(
	user *user.User,
	name *names.Name) ([]*decision.Recommendation, error) {

	recs := []*decision.Recommendation{}
	_, err := datastore.NewQuery(entityTypeRecommendations).
		Filter(filterRecommendationUser, user.Email).
		Filter(filterNameEquals, name.Name).
		Filter(filterGenderEquals, name.Gender.GoString()).
		GetAll(mgr.ctx, &recs)

	if err != nil {
		return nil, err
	}

	return recs, nil

}
