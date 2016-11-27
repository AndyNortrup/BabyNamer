package main

import (
	"os"
	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
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
	q := datastore.NewQuery(EntityTypeNameDetails).KeysOnly()
	details := []NameDetails{}
	keys, err := q.GetAll(ctx, details)
	if err != nil {
		return err
	}

	err = datastore.DeleteMulti(ctx, keys)
	return err
}
