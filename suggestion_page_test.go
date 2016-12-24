package babynamer

import (
	"errors"
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/recommendation"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/user"
	"testing"
)

func TestSuggestionPage_RandomPage(t *testing.T) {
	u := &user.User{Email: "test@test.com"}
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	mpm := &mockPersistenceManager{randCount: 0}
	sp := NewSuggestionPage(u, mpm, ctx)
	for x := 0; x < len(randomNames); x++ {
		rndName := sp.randomName()

		if rndName != randomNames[x] {
			t.Fail()
			t.Logf("Wrong random name returned. \tExpected: %v \tRecieved: %v", randomNames[x], rndName)
		}
	}
}

var randomNames []*names.Name = []*names.Name{{Name: "TestName1"}, nil}
var randomErrors []error = []error{nil, errors.New("Random name error.")}

type mockPersistenceManager struct {
	randCount int
}

func (mpm *mockPersistenceManager) GetName(name string, gender names.Gender) ([]*names.Name, error) {
	return nil, nil
}
func (mpm *mockPersistenceManager) GetRandomName(names.Gender) (*names.Name, error) {
	resultName := randomNames[mpm.randCount]
	resultError := randomErrors[mpm.randCount]

	mpm.randCount++

	return resultName, resultError
}

func (mpm *mockPersistenceManager) AddName(name *names.Name) error {
	return nil
}
func (mpm *mockPersistenceManager) GetRecommendedNames() []*names.Name {
	return nil
}
func (mpm *mockPersistenceManager) UpdateDecision(names.Name, decision.Decision) {}