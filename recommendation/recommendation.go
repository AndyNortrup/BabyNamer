package decision

import (
	"github.com/AndyNortrup/baby-namer/names"
	"google.golang.org/appengine/user"
)

type Recommendation struct {
	Email       string
	Recommended bool
}

func NewRecommendation(user *user.User, recommended bool) *Recommendation {
	return &Recommendation{Email: user.Email, Recommended: recommended}
}

type RecommendationList struct {
	Name      names.Name
	Decisions []*Recommendation
}
