package babynamer

import "testing"

var recommender = "Recommender"
var approver = "Approver"

func TestDecision_Approve(t *testing.T) {
	blank := &Decision{}
	blank.Approve(recommender)

	if blank.RecommendedBy != recommender {
		t.Logf("Expected: %v\t Recieved:%v", recommender, blank.RecommendedBy)
		t.Fail()
	}

	blank.Approve(approver)
	if blank.ApprovedBy != approver {
		t.Logf("Expected: %v\t Recieved:%v", approver, blank.ApprovedBy)
		t.Fail()
	}
}

func TestDecision_Reject(t *testing.T) {
	blank := &Decision{}
	blank.Reject(recommender)

	if blank.RejectedBy != recommender {
		t.Logf("Expected: %v\t Recieved: %v", approver, blank.RejectedBy)
		t.Fail()
	}

	if !blank.IsRejected() {
		t.Logf("Expected: true \t Recieved: %v", blank.IsRejected())
	}
}
