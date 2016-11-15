package main

import (
	"net/http"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func TestGetRecommendedName(t *testing.T) {

	userEmail := "me@example.com"

	req, err := inst.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Unable to create request for GetRecommendedName: %v", err)
	}
	ctx := appengine.NewContext(req)

	gen := &NameGenerator{
		ctx:    ctx,
		user:   userEmail,
		client: urlfetch.Client(ctx),
	}

	desired := &NameDetails{
		Name:          "test",
		RecommendedBy: "user@test.com",
	}

	gen.addNameToStore(desired)

	rejected := &NameDetails{
		Name:       "reject",
		RejectedBy: "user@test.com",
	}

	gen.addNameToStore(rejected)

	undesired := &NameDetails{
		Name:          "banannas",
		RecommendedBy: userEmail,
	}
	//put a test value into the store

	err = gen.addNameToStore(undesired)

	if err != nil {
		t.Fatalf("Failed to add name to store: %v", err)
	}

	//test that we get that value back
	ndResult, err := gen.getRecommendedName("test")

	if err != nil {
		t.Fatalf("Error getting recommended name: %v", err)
	}
	if ndResult == nil {
		t.Fatalf("nil result returned.")
	}
	if ndResult.Name != desired.Name {
		t.Logf("Incorrect name returned. Expected: %v Recieved: %v",
			desired.Name, desired.Name)
		t.Fail()
	}

	if err != nil {
		t.Fatal(err)
	}
}
