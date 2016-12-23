package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
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
