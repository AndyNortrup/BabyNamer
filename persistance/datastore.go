package persist

import (
	"errors"
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"math/rand"
	"time"
)

type DatastorePersistenceManager struct {
	ctx context.Context
}

const entityTypeName string = "Name"
const entityTypeStats string = "Stats"
const filterNameEquals string = "Name = "
const filterGenderEquals string = "Gender ="
const filterRandomNumber string = "Random >="

func NewDatastoreManager(ctx context.Context) *DatastorePersistenceManager {
	return &DatastorePersistenceManager{
		ctx: ctx,
	}
}

func (mgr *DatastorePersistenceManager) GetName(name string, gender names.Gender) ([]*names.Name, error) {
	query := datastore.NewQuery(entityTypeName).
		Filter(filterNameEquals, name).
		Filter(filterGenderEquals, gender)
	return mgr.executeGetNameQuery(mgr.ctx, query)
}

func (mgr *DatastorePersistenceManager) GetRandomName(gender names.Gender) (*names.Name, error) {
	query := datastore.NewQuery(entityTypeName).
		Filter(filterGenderEquals, gender).
		Filter(filterRandomNumber, randomFloat()).
		Limit(1)

	results, err := mgr.executeGetNameQuery(mgr.ctx, query)
	return results[0], err
}

func (mgr *DatastorePersistenceManager) AddName(name *names.Name) error {
	dName := newDatastoreName(name)
	key := datastore.NewIncompleteKey(mgr.ctx, name.Key(), nil)
	key, err := datastore.Put(mgr.ctx, key, dName)
	if err != nil {
		log.Errorf(mgr.ctx, "Failed to put name %v into datastore: %v", name.Name, err)
		return err
	}

	err = mgr.addStatsToDatastore(mgr.ctx, key, name.Stats)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *DatastorePersistenceManager) GetRecommendedNames() []*names.Name {
	return nil
}

func (mgr *DatastorePersistenceManager) UpdateDecision(name names.Name, decision decision.Decision) {

}

func (mgr *DatastorePersistenceManager) executeGetNameQuery(ctx context.Context, query *datastore.Query) ([]*names.Name, error) {
	results := []*datastoreName{}
	keys, err := query.GetAll(mgr.ctx, &results)

	if err != nil {
		return nil, errors.New("Failed to get name: " + err.Error())
	}

	names := []*names.Name{}
	for index, name := range results {
		name := &name.Name
		mgr.getStatsForKey(name, keys[index])
		names = append(names, name)
	}

	return names, nil
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

func (mgr *DatastorePersistenceManager) addStatsToDatastore(ctx context.Context, parent *datastore.Key, stats map[int]*names.Stat) error {

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
