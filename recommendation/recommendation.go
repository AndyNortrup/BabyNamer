package decision

import (
	"github.com/AndyNortrup/baby-namer/names"
	"google.golang.org/appengine/user"
)

type Recommendation struct {
	Name        string
	Gender      string
	Email       string
	Recommended bool
}

func NewRecommendation(user *user.User, name *names.Name, recommended bool) *Recommendation {
	return &Recommendation{
		Email:       user.Email,
		Recommended: recommended,
		Name:        name.Name,
		Gender:      name.Gender.GoString(),
	}
}
