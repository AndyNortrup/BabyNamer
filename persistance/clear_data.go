package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func ClearDatastore(ctx context.Context) error {
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

	return nds.DeleteMulti(ctx, keys)
}

func deleteStatsData(ctx context.Context) error {
	q := datastore.NewQuery(entityTypeStats).KeysOnly()
	stats := []names.Stat{}
	key, err := q.GetAll(ctx, stats)
	if err != nil {
		return err
	}

	return nds.DeleteMulti(ctx, key)
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
