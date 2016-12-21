package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"testing"
)

var input = []string{"Mary::F", "Anna::F", "Emma::F", "Pat::M", "Pat::F"}
var testNames = []string{"Mary", "Anna", "Emma", "Pat", "Pat"}
var first = []int{1880, 1880, 1881, 1880, 1881}
var last = []int{1881, 1882, 1882, 1882, 1881}
var max = []int{1880, 1882, 1882, 1882, 1881}
var gender = []string{"F", "F", "F", "M", "F"}
var genderFilter = []names.Gender{
	names.FemaleFilter,
	names.FemaleFilter,
	names.FemaleFilter,
	names.MaleFilter,
	names.FemaleFilter,
}

func TestGetName(t *testing.T) {
	ctx := newTestContext()
	mgr := NewDatastoreManager(ctx)
	setupDatastoreTest(mgr, t)

	for index, name := range testNames {
		result, err := mgr.GetName(name, genderFilter[index])
		if err != nil {
			t.Logf("Failed to get name: %v - %v", name, err)
			t.FailNow()
		}

		for x, value := range result {
			checkResults(index, input[x], value, t)
		}
	}

	deleteAllNameDetails(ctx)
}

func TestDatastorePersistenceManager_GetRandomName(t *testing.T) {
	ctx := newTestContext()
	mgr := NewDatastoreManager(ctx)
	setupDatastoreTest(mgr, t)

	//Check that we get random values back
	_, err := mgr.GetRandomName(names.FemaleFilter)
	if err != nil {
		t.Logf("Failed to get random name: %v", err)
		t.Fail()
	}

}

func setupDatastoreTest(mgr DataManager, t *testing.T) {
	inputData := LoadNames()
	for _, name := range inputData {
		err := mgr.AddName(name)
		if err != nil {
			t.Logf("Failed to write data to datastore: %v", err)
			t.FailNow()
		}
	}
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

	if result.MostPopularYear().Year != max[idx] {
		t.Logf("%v: Expected most popular year: %v Recieved: %v", name, last[idx], result.MostPopularYear().Year)
		t.Fail()
	}

	if result.Gender != gender[idx] {
		t.Logf("%v: Expected gender: %v, Recieved: %v", name, gender[idx], result.Gender)
	}
}
