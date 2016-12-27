package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"testing"
	"time"
)

var input = []string{"Mary::F", "Anna::F", "Emma::F", "Pat::M", "Pat::F"}
var testNames = []string{"Mary", "Anna", "Emma", "Pat", "Pat"}
var first = []int{1880, 1880, 1881, 1880, 1881}
var last = []int{1881, 1882, 1882, 1882, 1881}
var max = []int{1880, 1882, 1882, 1882, 1881}
var gender = []names.Gender{names.Female, names.Female, names.Female, names.Male, names.Female}
var genderFilter = []names.Gender{
	names.Female,
	names.Female,
	names.Female,
	names.Male,
	names.Female,
}

func TestGetName(t *testing.T) {
	ctx := newTestContext()
	mgr := NewDatastoreManager(ctx)
	LoadNames(ctx)

	for index, name := range testNames {
		result, err := mgr.GetName(name, genderFilter[index])
		if err != nil {
			t.Logf("Failed to get name: %v - %v", name, err)
			t.FailNow()
		}
		checkResults(index, input[index], result, t)
	}
	deleteAllNameDetails(ctx)
}

func TestDatastorePersistenceManager_GetRandomName(t *testing.T) {
	ctx := newTestContext()
	mgr := NewDatastoreManager(ctx)
	LoadNames(ctx)

	//Check that we get random values back
	// Because we only have 5 names in the test dataset, we need to try a few times to make sure we get one.
	received := false

	for x := 0; x < 5; x++ {
		_, err := mgr.GetRandomName(names.Female)
		if err == nil {
			received = true
			break
		}
		if err != nil && err != NoRandomName {
			t.Logf("Failed to get random name: %v", err)
			t.Fail()
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !received {
		t.Fail()
		t.Log("No random name returned.")
	}

	deleteAllNameDetails(ctx)
}

func checkResults(idx int, name string, result *names.Name, t *testing.T) {
	if result.Name != testNames[idx] {
		t.FailNow()
	}
	if result.FirstYear().Year != first[idx] {
		t.Logf("Expected first year: %v Recieved: %v", first[idx], result.FirstYear().Year)
		t.Fail()
	}

	if result.LatestYear().Year != last[idx] {
		t.Logf("%v: Expected latest year: %v Recieved: %v", name, last[idx], result.LatestYear().Year)
		t.Fail()
	}

	if result.MostOccurrences().Year != max[idx] {
		t.Logf("%v: Expected most popular year: %v Recieved: %v", name, last[idx], result.MostOccurrences().Year)
		t.Fail()
	}

	if result.Gender != gender[idx] {
		t.Logf("%v: Expected gender: %v, Recieved: %v", name, gender[idx], result.Gender)
	}
}

func TestDatastorePersistenceManager_UpdateDecision(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatalf("action=TestDatastorePersistenceManager_UpdateDecision error=%v", err)
	}
	defer done()

	usr := &user.User{Email: "test1@test.com"}

	data := NewDatastoreManager(ctx)
	name := names.NewName("Mary", names.Female)
	rec := decision.NewRecommendation(usr, true)

	data.AddName(name)

	err = data.UpdateDecision(name, rec)
	if err != nil {
		t.Fatalf("action=TestUpdateDecision error=%v", err)
	}

	out, err := getAllRecommendations(ctx)
	if err != nil {
		t.Fatalf("action=get_all_recommendations error=%v", err)
	}

	if len(out) == 0 {
		t.Logf("action=TestDatastorePersistenceManager_UpdateDecision no recommendations retrived "+"Recieved=%v", len(out))
		t.FailNow()
	}

	if out[0].Recommended != true {
		t.Logf("action=TestDatastorePersistenceManager_UpdateDecision "+"attribute=Recommended expected=%v recieved=%v",
			true, out[0].Recommended)
		t.FailNow()
	}

	rec.Recommended = false
	data.UpdateDecision(name, rec)
	out, err = getAllRecommendations(ctx)
	if err != nil {
		t.Fatalf("action=TestDatastorePersistenceManager_UpdateDecision error=%v", err)
	}

	if len(out) != 1 {
		t.Log("action=TestDatastorePersistenceManager_UpdateDecision no recommendations retrived")
		t.FailNow()
	}

	if out[0].Recommended != false {
		t.Log("action=TestDatastorePersistenceManager_UpdateDecision attribute=Recommended expected=%v recieved=%v",
			false, out[0].Recommended)
		t.FailNow()
	}

	deleteAllNameDetails(ctx)
}

func getAllRecommendations(ctx context.Context) ([]*decision.Recommendation, error) {
	q := datastore.NewQuery(entityTypeRecommendations)
	out := []*decision.Recommendation{}
	_, err := q.GetAll(ctx, &out)
	return out, err
}

func TestDatastorePersistenceManager_GetRecommendedNames(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatalf("action=TestGetRecommendedNames error=%v", err)
	}
	defer done()

	recommendedName := names.NewName("Recommended", names.Female)
	rejectedName := names.NewName("Rejected", names.Female)
	undecidedName := names.NewName("Undecided", names.Female)

	user1 := &user.User{Email: "user1@test.com"}
	user2 := &user.User{Email: "user2@test.com"}

	rec1 := decision.NewRecommendation(user1, true)
	rec2 := decision.NewRecommendation(user1, false)

	mgr := NewDatastoreManager(ctx)

	mgr.AddName(recommendedName)
	mgr.AddName(rejectedName)
	mgr.AddName(undecidedName)

	mgr.UpdateDecision(recommendedName, rec1)
	mgr.UpdateDecision(recommendedName, rec2)

	recommendNames, err := mgr.GetRecommendedNames(user1, user2)

	if len(recommendNames) != 1 {
		t.FailNow()
		t.Logf("action=TestDatastorePersistenceManager_GetRecommendedNames "+"failure=incorrect number of recommended names returned "+"expected=1 recieved=%v", len(recommendNames))
	}

	if recommendNames[0] != recommendedName {
		t.FailNow()
		t.Logf("action=TestDatastorePersistenceManager_GetRecommendedNames "+"failure=wrong name returned "+"expected=%v "+"recieved=%v ", recommendedName, recommendNames[0])
	}

	deleteAllNameDetails(ctx)
}
