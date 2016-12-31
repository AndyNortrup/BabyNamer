package persist

import (
	"os"
	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"

	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"net/http"
)

var inst aetest.Instance

func TestMain(m *testing.M) {
	var err error
	inst, err = aetest.NewInstance(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		os.Exit(2)
	}
	defer tearDown()

	m.Run()
}

func tearDown() {
	if inst != nil {
		inst.Close()
	}
}

func clearDatastore(ctx context.Context) error {
	err := deleteNameData(ctx)
	if err != nil {
		return err
	}
	err = deleteStatsData(ctx)
	if err != nil {
		return err
	}
	err = deleteRecommendationData(ctx)
	if err != nil {
		return err
	}
	return nil
}

func deleteNameData(ctx context.Context) error {
	q := datastore.NewQuery(entityTypeName).KeysOnly()
	details := []names.Name{}
	keys, err := q.GetAll(ctx, details)
	if err != nil {
		return err
	}

	return datastore.DeleteMulti(ctx, keys)
}

func deleteStatsData(ctx context.Context) error {
	q := datastore.NewQuery(entityTypeStats).KeysOnly()
	stats := []names.Stat{}
	key, err := q.GetAll(ctx, stats)
	if err != nil {
		return err
	}

	return datastore.DeleteMulti(ctx, key)
}

func deleteRecommendationData(ctx context.Context) error {
	q := datastore.NewQuery(entityTypeRecommendations).KeysOnly()
	stats := []decision.Recommendation{}
	key, err := q.GetAll(ctx, stats)
	if err != nil {
		return err
	}

	return datastore.DeleteMulti(ctx, key)
}

func newTestContext() context.Context {
	req, _ := inst.NewRequest(http.MethodGet, "/", nil)
	return appengine.NewContext(req)
}
