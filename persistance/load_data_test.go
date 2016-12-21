package persist

import "testing"

func TestLoadNames(t *testing.T) {
	result := LoadNames()
	if len(result) == 0 {
		t.Log("No names returned.")
		t.Fail()
	}

	if len(result) != 5 {
		t.Logf("Wrong number of names returned. Expected 5\t Recieved: %v", len(result))
	}
}
