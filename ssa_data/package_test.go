package ssa_data

import (
	"context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"net/http"
	"os"
	"testing"
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

func deleteAllNameDetails(ctx context.Context) error {
	q := datastore.NewQuery(entityTypeName).KeysOnly()
	details := []Name{}
	keys, err := q.GetAll(ctx, details)
	if err != nil {
		return err
	}

	err = datastore.DeleteMulti(ctx, keys)
	return err
}

func newTestContext() context.Context {
	req, _ := inst.NewRequest(http.MethodGet, "/", nil)
	return appengine.NewContext(req)
}
