package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
)

type DataManager interface {
	GetName(name string, gender names.Gender) ([]*names.Name, error)
	GetRandomName(names.Gender) (*names.Name, error)
	AddName(name *names.Name) error
	GetRecommendedNames() []*names.Name
	UpdateDecision(names.Name, decision.Decision)
}
