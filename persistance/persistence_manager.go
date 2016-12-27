package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"google.golang.org/appengine/user"
)

type DataManager interface {
	GetName(name string, gender names.Gender) (*names.Name, error)
	GetRandomName(names.Gender) (*names.Name, error)
	AddName(name *names.Name) error
	GetRecommendedNames(user, partner *user.User) ([]*names.Name, error)
	UpdateDecision(*names.Name, *decision.Recommendation) error
}
