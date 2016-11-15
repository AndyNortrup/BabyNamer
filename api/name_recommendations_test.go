package main

import "testing"

//TestDuplicateDecisions tests to see if the checkDuplicateDecisions method
// accurately detects a duplicate decision by a user.  Users should not be able to
// both recommend and approve or recommend and reject a single name.
func TestDuplicateDecisions(t *testing.T) {
	user := "user"
	otherUser := "otherUser"
	names := []*NameDetails{
		{Name: "Approved", ApprovedBy: user},
		{Name: "Rejected", RejectedBy: user},
		{Name: "RecommendedBy", RecommendedBy: user},
		{Name: "Good", RecommendedBy: otherUser},
	}

	results := []bool{true, true, true, false}

	for index, name := range names {
		result := checkDuplicateDecision(name, user)
		if (result != nil) != results[index] {
			t.Errorf("checkDuplicateDecision failed to identify duplicate decision. \n\t%v", result)
			t.Fail()
		}
	}
}
