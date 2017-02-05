package babynamer

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/persistance"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
)

func newDataManager(ctx context.Context) DataManager {
	return persist.NewDatastoreManager(ctx)
}

type DataManager interface {
	//GetName retrieves a specific name and all of it's stats from the store
	GetName(name string, gender names.Gender) (*names.Name, error)

	//GetRandomName gets a random name from the datastore, this should not take into account
	// recommendations by the user or their partner
	GetRandomName(names.Gender) (*names.Name, error)

	//AddName adds the name to the store
	AddName(name *names.Name) error

	//GetNameRecommendations gets recommendations by the user for a given name.
	GetNameRecommendations(user *user.User, name *names.Name) ([]*decision.Recommendation, error)

	//GetPartnerRecommendedNames returns a list of names that have been recommended by the users partner.
	GetPartnerRecommendedNames(user, partner *user.User) ([]*names.Name, error)

	//UpdateDecision get public
	UpdateDecision(*decision.Recommendation) error

	//GetShortList returns a list of all of the names that the user and partner have both agreed too.
	GetShortList(user, partner *user.User) (names.NameList, error)
}
