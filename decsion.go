package babynamer

type Decision struct {
	RecommendedBy string
	ApprovedBy    string
	RejectedBy    string
}

func (d *Decision) IsRecommended() bool {
	return d.RecommendedBy > ""
}

func (d *Decision) Approve(approver string) {
	if d.RecommendedBy == "" {
		d.RecommendedBy = approver
	} else {
		d.ApprovedBy = approver
	}
}

func (d *Decision) IsApproved() bool {
	return d.ApprovedBy > ""
}

func (d *Decision) Reject(rejector string) {
	d.RejectedBy = rejector
}

func (d *Decision) IsRejected() bool {
	return d.RejectedBy > ""
}
