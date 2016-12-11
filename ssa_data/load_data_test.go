package ssa_data

import (
	"testing"
)

func TestLoadNames(t *testing.T) {
	result := LoadNames()
	if len(result) == 0 {
		t.Log("No names returned.")
		t.Fail()
	}

	input := []string{"Mary", "Anna", "Emma"}
	first := []int{1880, 1880, 1881}
	last := []int{1881, 1882, 1882}
	max := []int{1880, 1882, 1882}

	for idx, name := range input {
		if result[name].Name != input[idx] {
			t.Logf("Failed to retireve name: %v", name)
			t.FailNow()
		}
		if result[name].FirstYear().Year != first[idx] {
			t.Logf("Expected first year: %v Recieved: %v", first[idx], result[name].FirstYear().Year)
			t.Fail()
		}

		if result[name].LatestYear().Year != last[idx] {
			t.Logf("%v: Expected latest year: %v Recieved: %v", name, last[idx], result[name].LatestYear().Year)
			t.Fail()
		}

		if result[name].MostPopularYear().Year != max[idx] {
			t.Logf("%v: Expected most popular year: %v Recieved: %v", name, last[idx], result[name].MostPopularYear().Year)
			t.Fail()
		}
	}
}
