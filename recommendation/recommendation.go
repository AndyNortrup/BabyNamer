package decision

import (
	"github.com/AndyNortrup/baby-namer/names"
	"google.golang.org/appengine/user"
)

type Decision struct {
	Name          names.Name `datastore: "-"`
	RecommendedBy []*user.User
	RejectedBy    []*user.User
}

func NewDecision() *Decision {
	return &Decision{}
}

func (rec *Decision) Recommend(user *user.User) {
	rec.RecommendedBy = append(rec.RecommendedBy, user)
}

func (rec *Decision) Reject(user *user.User) {
	rec.RejectedBy = append(rec.RejectedBy, user)
}
