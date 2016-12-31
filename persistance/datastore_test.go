package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"golang.org/x/net/context"
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
	LoadNames(ctx, "persistance/names")

	for index, name := range testNames {
		result, err := mgr.GetName(name, genderFilter[index])
		if err != nil {
			t.Logf("Failed to get name: %v - %v", name, err)
			t.FailNow()
		}
		checkResults(index, input[index], result, t)
	}
	clearDatastore(ctx)
}

func TestDatastorePersistenceManager_GetRandomName(t *testing.T) {
	ctx := newTestContext()
	mgr := NewDatastoreManager(ctx)
	LoadNames(ctx, "persistance/names")

	//Check that we get random values back
	// Because we only have 5 names in the test data set, we need to try a few times to make sure we get one.
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

	clearDatastore(ctx)
}

func checkResults(idx int, name string, result *names.Name, t *testing.T) {
	if result.Name != testNames[idx] {
		t.FailNow()
	}
	if len(result.Stats) == 0 {
		t.Log("No stats returned.")
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
	ctx := newTestContext()

	usr := &user.User{Email: "test1@test.com"}

	data := NewDatastoreManager(ctx)
	name := names.NewName("Mary", names.Female)
	rec := decision.NewRecommendation(usr, true)

	err := data.AddName(name)
	if err != nil {
		t.Fatalf("action=TestDatastorePersistenceManager_UpdateDecision error=%v", err)
	}

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
		t.Logf("action=TestDatastorePersistenceManager_UpdateDecision attribute=Recommended expected=%v recieved=%v",
			false, out[0].Recommended)
		t.FailNow()
	}

	clearDatastore(ctx)
}

func getAllRecommendations(ctx context.Context) ([]*decision.Recommendation, error) {
	q := datastore.NewQuery(entityTypeRecommendations)
	out := []*decision.Recommendation{}
	_, err := q.GetAll(ctx, &out)
	return out, err
}

func TestDataStorePersistenceManager_GetUserRecommendations(t *testing.T) {
	ctx := newTestContext()

	recommendedName := names.NewName("Recommended", names.Female)

	user1 := &user.User{Email: "user1@test.com"}
	rec1 := decision.NewRecommendation(user1, true)

	mgr := NewDatastoreManager(ctx)
	err := mgr.AddName(recommendedName)
	if err != nil {
		t.Fatalf("action=TestDataStorePersistenceManager_GetUserRecommendations error=%v", err)
	}
	mgr.UpdateDecision(recommendedName, rec1)

	recs, _, err := mgr.getUserRecommendations(mgr.getRecommendationQuery(user1, true))
	if err != nil {
		t.Fatalf("action=TestDataStorePersistenceManager_GetUserRecommendations "+"error=%v", err)
	}

	if len(recs) != 1 {
		t.Logf("action=TestDataStorePersistenceManager_GetUserRecommendations "+"incorrect number of recommendations expected=1"+"recieved=%v", len(recs))
		t.FailNow()
	}

	rec1.Recommended = false
	mgr.UpdateDecision(recommendedName, rec1)
	recs, _, err = mgr.getUserRecommendations(mgr.getRecommendationQuery(user1, true))
	if len(recs) != 0 {
		t.Logf("action=TestDataStorePersistenceManager_GetUserRecommendations "+"incorrect number of recommendations expected=0"+"recieved=%v", len(recs))
		t.FailNow()
	}

	clearDatastore(ctx)
	ctx.Done()
}

func TestDatastorePersistenceManager_GetRecommendedNames(t *testing.T) {
	ctx := newTestContext()

	testNames := []*names.Name{
		names.NewName("Recommended", names.Female),
		names.NewName("Rejected", names.Female),
		names.NewName("Undecided", names.Female),
	}

	usr := &user.User{Email: "usr@test.com"}
	partner := &user.User{Email: "partner@test.com"}

	rec1 := decision.NewRecommendation(partner, true)
	rec2 := decision.NewRecommendation(partner, false)

	mgr := NewDatastoreManager(ctx)

	for _, name := range testNames {
		err := mgr.AddName(name)
		if err != nil {
			t.Fatal(err)
		}
	}

	mgr.AddName(testNames[0])
	mgr.AddName(testNames[1])
	mgr.AddName(testNames[2])

	mgr.UpdateDecision(testNames[0], rec1)
	mgr.UpdateDecision(testNames[2], rec2)

	recommendedNames, err := mgr.GetRecommendedNames(usr, partner)
	if err != nil {
		t.Fatal(err)
	}

	if len(recommendedNames) != 1 {
		t.Logf("action=TestDatastorePersistenceManager_GetRecommendedNames "+"failure=incorrect number of recommended testNames returned "+"expected=1 recieved=%v", len(recommendedNames))
		t.FailNow()
	}

	if !compareNames(testNames[0], recommendedNames[0]) {
		t.Logf("action=TestDatastorePersistenceManager_GetRecommendedNames "+"failure=wrong name returned "+"expected=%v "+"recieved=%v ", testNames[0], recommendedNames[0])
		t.FailNow()
	}

	clearDatastore(ctx)
	ctx.Done()
}

func TestDatastorePersistenceManager_addStat(t *testing.T) {

	ctx := newTestContext()

	nIn := names.NewName("Mary", names.Female)
	nIn.AddStat(names.NewNameStat(1, 1, 1))

	mgr := NewDatastoreManager(ctx)
	err := mgr.AddName(nIn)

	nOut, err := mgr.GetName(nIn.Name, nIn.Gender)
	if err != nil {
		t.Fatalf("action=TestDatastorePersistenceManager_addStat error=%v", err)
	}

	if len(nOut.Stats) != 1 {
		t.Logf("action=TestDatastorePersistenceManager_addStat wrong number of stats returned "+"expected=1 recieved=%v", len(nOut.Stats))
		t.FailNow()
	}

	if *nOut.Stats[1] != *nIn.Stats[1] {
		t.Log("action=TestDatastorePersistenceManager_addStat stat not properly recorded")
		t.FailNow()
	}

	nIn.Stats[1] = names.NewNameStat(1, 2, 2)
	err = mgr.addStatToDatastore(nIn, nIn.Stats[1])

	if err != nil {
		t.Fatal(err)
	}

	nOut, err = mgr.GetName(nIn.Name, nIn.Gender)

	if len(nOut.Stats) != 1 {
		t.Log("action=TestDatastorePersistenceManager_addStat wrong number of stats returned")
		t.FailNow()
	}

	if *nOut.Stats[1] != *nIn.Stats[1] {
		t.Logf("action=TestDatastorePersistenceManager_addStat stat not properly recorded "+"\nexpected: %v"+"\nrecieved: %v",
			nIn.Stats[1], nOut.Stats[1])
		t.FailNow()
	}

	clearDatastore(ctx)
}
