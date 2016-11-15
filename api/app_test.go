package main

import (
	"os"
	"testing"

	"google.golang.org/appengine/aetest"
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
