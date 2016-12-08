package babynamer

import (
	"net/http"
	"testing"

	"google.golang.org/appengine"
)

func TestGetName(t *testing.T) {

	req, err := inst.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Unable to create request for GetName: %v", err)
	}
	ctx := appengine.NewContext(req)
	username := "username"

	mgr := NewDatastoreNameManager(ctx, username)

	name := &NameDetails{
		Name: "Test",
	}

	err = mgr.addNameToStore(name)
	if err != nil {
		t.Fatalf("Failed to add name to store: %v", err)
	}

	result, _, err := mgr.getNameFromDatastore("Test")
	if err != nil {
		t.Fatalf("Failed to retrive name from datastore: %v", err)
	}

	if result.Name != name.Name {
		t.Logf("Retrived name does not match. \n\tExpected:%v \t\nRecieved: %v",
			name.Name,
			result.Name)
	}

}
