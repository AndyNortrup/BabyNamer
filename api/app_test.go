package main

import (
	"net/http"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func TestGetRecommendedName(t *testing.T) {

	userEmail := "me@example.com"

	// user := &user.User{Email: userEmail}
	// aetest.Login(user, req)
	//
	req := newTestRequest(t)

	ctx := appengine.NewContext(req)

	t.Logf("State of ctx: %#v", ctx)

	gen := NewNameGenerator(ctx)

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

	err := gen.addNameToStore(undesired)

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

func newTestRequest(t *testing.T) *http.Request {
	inst, err := aetest.NewInstance(
		&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer inst.Close()

	req, err := inst.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	return req
}
