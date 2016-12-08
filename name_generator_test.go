package babynamer

import (
	"github.com/AndyNortrup/baby-namer/usage"
	"golang.org/x/net/context"
	"testing"
)

var names = []*NameDetails{
	{Name: "Recommended",
		Decision: Decision{
			RecommendedBy: "user1",
		},
		Usages: []usage.Usage{desiredUsage},
	},
	{Name: "Rejected",
		Decision: Decision{
			RecommendedBy: "user1",
		},
		Usages: []usage.Usage{unwantedUsage}},
	{Name: "Undecided Correct Usage", Usages: []usage.Usage{desiredUsage}},
	{Name: "Undecided Wrong Usage", Usages: []usage.Usage{unwantedUsage}},
}

var desiredUsage = usage.Usage{UsageFull: "desiredUsage"}
var unwantedUsage = usage.Usage{UsageFull: "unwanted"}

func TestGetRecommendedName(t *testing.T) {

	userEmail := "me@example.com"
	otherEmail := "user@test.com"
	desiredName := "recommendedByPartner"

	ctx := newTestContext()

	gen := &NameGenerator{
		ctx:  ctx,
		user: userEmail,
	}

	mgr := &DatastoreNameManager{
		ctx:      ctx,
		username: userEmail,
	}

	testCases := []*NameDetails{
		{
			Name: desiredName,
			Decision: Decision{
				RecommendedBy: otherEmail,
			},
		},
		{
			Name: "reject",
			Decision: Decision{
				RejectedBy: otherEmail,
			},
		},
		{
			Name: "recommendedBySelf",
			Decision: Decision{
				RecommendedBy: userEmail,
			},
		},
	}

	for _, name := range testCases {
		err := mgr.addNameToStore(name)

		if err != nil {
			t.Fatalf("Failed to add name to store: %v", err)
		}

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

	cleanupDataStore(ctx, t)

}

func TestGetUndecdiedName(t *testing.T) {
	ctx := newTestContext()
	username := "user1"

	addUndecidedNameData(username, ctx, t)

	gen := NewNameGenerator(ctx, username)
	result, err := gen.getUndecidedName(&desiredUsage)

	if err != nil {
		t.Fatal(err)
	}

	evaluateUndecidedNameResults(result, t)

	cleanupDataStore(ctx, t)

}

func addUndecidedNameData(username string, ctx context.Context, t *testing.T) {

	mgr := NewDatastoreNameManager(ctx, username)

	for _, name := range names {
		err := mgr.addNameToStore(name)
		if err != nil {
			deleteAllNameDetails(ctx)
			t.Fatalf("Failed to add names to datastore: %v", err)
		}
	}
}

func evaluateUndecidedNameResults(result *NameDetails, t *testing.T) {

	if result == nil {
		t.Log("No name returned.")
		t.FailNow()
	}

	if result.Name != names[2].Name {
		t.Logf("Expected: %v\nRecieved: %v", result, names[2])
		t.Fail()
	}
}

func cleanupDataStore(ctx context.Context, t *testing.T) {
	err := deleteAllNameDetails(ctx)
	if err != nil {
		t.Fatalf("Failed to delete name details: %v", err)
	}
}
