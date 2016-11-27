package main

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"testing"
)

func TestGetShortList(t *testing.T) {

	req, err := inst.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Unable to create request for GetName: %v", err)
	}
	ctx := appengine.NewContext(req)

	names := []NameDetails{
		{Name: "Alpha", RecommendedBy: "one", ApprovedBy: "Two"},
		{Name: "Bravo", RecommendedBy: "one", RejectedBy: "Two"},
		{Name: "Charlie", RecommendedBy: "one", ApprovedBy: "Two"},
	}

	for _, name := range names {
		key := datastore.NewIncompleteKey(ctx, EntityTypeNameDetails, nil)
		key, err := datastore.Put(ctx, key, &name)

		if err != nil {
			t.Fatalf("Failed to add test data to datastore: %v", err)
		}
	}

	shortList, err := getShortList(ctx)

	if err != nil {
		t.Log("Failed to get ShortList from Datastore.")
		t.Fail()
	}

	if shortList.Len() != 2 {
		t.Log("Incorrect number of items in ShortList")
		t.FailNow()
	}

	if shortList[0].Name != names[0].Name {
		t.Logf("Wrong name returned at top of ShortList.  \n\tExpected: %v \n\tRecieved: %v",
			names[0].Name, shortList[0].Name)
	}

	deleteAllNameDetails(ctx)
}
