package settings_test

import (
	"os"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
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

func newTestContext() context.Context {
	req, _ := inst.NewRequest(http.MethodGet, "/", nil)
	return appengine.NewContext(req)
}
