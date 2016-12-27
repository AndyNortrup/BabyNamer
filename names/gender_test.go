package names

import "testing"

func TestGetGender(t *testing.T) {
	inputs := []string{"M", "F"}
	outputs := []Gender{Male, Female}

	for key, input := range inputs {
		if GetGender(input) != outputs[key] {
			t.Logf("Failed to identify gender.  Expected=%v\t Recieved=%v",
				outputs[0], GetGender(input))
			t.Fail()
		}
	}
}
