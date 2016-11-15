package main

import (
	"net/http"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func TestGetRecommendedName(t *testing.T) {

	userEmail := "me@example.com"
	otherEmail := "user@test.com"
	desiredName := "recommendedByPartner"

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

	testCases := []*NameDetails{
		{
			Name:          desiredName,
			RecommendedBy: otherEmail,
		},
		{
			Name:       "reject",
			RejectedBy: otherEmail,
		},
		{
			Name:          "recommendedBySelf",
			RecommendedBy: userEmail,
		},
	}

	for _, name := range testCases {
		gen.addNameToStore(name)
	}

	if err != nil {
		t.Fatalf("Failed to add name to store: %v", err)
	}

	//test that we get that value back
	ndResult, err := gen.getRecommendedName("test")

	if err != nil {
		t.Fatalf("Error getting recommended name: %v", err)
	}
	if ndResult == nil {
		t.Fatal("nil result returned.")
	}
	if ndResult.Name != desiredName {
		t.Logf("Incorrect name returned. Expected: %v Recieved: %v",
			desiredName, ndResult.Name)
		t.Fail()
	}

	if err != nil {
		t.Fatal(err)
	}
}
