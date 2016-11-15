package main

//checkDuplicateDecision prevents the same user for redering judgement on a name more than once.
func checkDuplicateDecision(details *NameDetails, username string) error {
	//Make sure that this user hasn't already recommended / approved
	// rejected this name.
	if details.RecommendedBy == username {
		return NewDuplicateDecisionError(username, "recommended")
	} else if details.ApprovedBy == username {
		return NewDuplicateDecisionError(username, "approved")
	} else if details.RejectedBy == username {
		return NewDuplicateDecisionError(username, "rejected")
	}
	return nil
}

type DuplicateDecisionError struct {
	message string
}

func NewDuplicateDecisionError(user, action string) *DuplicateDecisionError {
	return &DuplicateDecisionError{
		message: "User " + user + " already " + action + "this name",
	}
}

func (err DuplicateDecisionError) Error() string {
	return err.message
}
