package ssa_data

import (
	"testing"
)

var input = []string{"Mary::F", "Anna::F", "Emma::F", "Pat::M", "Pat::F"}
var names = []string{"Mary", "Anna", "Emma", "Pat", "Pat"}
var first = []int{1880, 1880, 1881, 1880, 1881}
var last = []int{1881, 1882, 1882, 1882, 1881}
var max = []int{1880, 1882, 1882, 1882, 1881}
var gender = []string{"F", "F", "F", "M", "F"}
var genderFilter = []Gender{FemaleFilter, FemaleFilter, FemaleFilter, MaleFilter, FemaleFilter}

func TestLoadNames(t *testing.T) {
	result := loadNames()
	if len(result) == 0 {
		t.Log("No names returned.")
		t.Fail()
	}

	for idx, name := range input {
		checkResults(idx, name, result[name], t)
	}
}

func TestGetName(t *testing.T) {
	ctx := newTestContext()
	inputData := loadNames()
	err := addNamesToDatastore(ctx, inputData)

	if err != nil {
		t.Logf("Failed to write data to datastore: %v", err)
		t.FailNow()
	}

	for index, name := range names {
		result, err := GetName(ctx, name, genderFilter[index])
		if err != nil {
			t.Logf("Failed to get name: %v - %v", name, err)
			t.FailNow()
		}

		for x, value := range result {
			checkResults(index, input[x], value, t)
		}

		//Check that we get random values back
		_, err = GetRandomName(ctx, FemaleFilter)
		if err != nil {
			t.Logf("Failed to get random name: %v", err)
			t.Fail()
		}
	}

	deleteAllNameDetails(ctx)
}

func checkResults(idx int, name string, result *Name, t *testing.T) {
	if result.Name != names[idx] {
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
