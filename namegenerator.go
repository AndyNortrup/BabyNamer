package main

import (
	"encoding/xml"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"
)

type NameGenerator struct {
	user   string
	ctx    context.Context
	client *http.Client
}

func NewNameGenerator(ctx context.Context) *NameGenerator {
	gen := &NameGenerator{ctx: ctx}
	gen.user = user.Current(ctx).String()
	gen.client = urlfetch.Client(ctx)

	return gen
}

func (gen *NameGenerator) getName(previous string) (*NameDetails, error) {
	name, err := gen.getRecommendedName(previous)
	if err != nil {
		return &NameDetails{}, err
	}

	if name != nil {
		log.Warningf(gen.ctx, "Returning recommended name: %v", name.Name)
		return name, nil
	}

	log.Infof(gen.ctx, "Returning random name")
	return gen.getNameFromService()
}

//getRecommendedName gets a list of names from the datastore that were
// recommended by a different user, and have not been rejected by any user
func (gen *NameGenerator) getRecommendedName(previous string) (*NameDetails, error) {
	query := datastore.NewQuery(NameDetailsEntityType).
		Filter("ApprovedBy =", "").
		Filter("RejectedBy =", "")

	for t := query.Run(gen.ctx); ; {
		details := &NameDetails{}
		_, err := t.Next(details)
		if err == datastore.Done {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		if details.RecommendedBy != gen.user && previous != details.Name {
			return details, nil
		}
	}
}

func (gen *NameGenerator) getNameFromService() (*NameDetails, error) {
	addr := "http://www.behindthename.com/api/random.php?" +
		"usage=" + gen.randomUsage() +
		"&key=an468794&number=1"
	log.Infof(gen.ctx, "Name request API: %v", addr)
	nameReq, err := gen.client.Get(addr)

	if err != nil {
		log.Warningf(gen.ctx, "Error getting random name: %v\n", err)
		return nil, err
	}

	defer nameReq.Body.Close()

	decoder := xml.NewDecoder(nameReq.Body)
	name := &BabyName{}
	err = decoder.Decode(name)

	return gen.getNameInfo(name)
}

func (gen *NameGenerator) getNameInfo(name *BabyName) (*NameDetails, error) {
	details := &NameDetails{}
	address := "http://www.behindthename.com/api/lookup.php?key=an468794&name=" +
		name.Name
	nameDetailsReq, err := gen.client.Get(address)

	if err != nil {
		log.Warningf(gen.ctx, "Unable to retrieve name details: %#v", err)
		return details, err
	}

	defer nameDetailsReq.Body.Close()
	decoder := xml.NewDecoder(nameDetailsReq.Body)

	err = decoder.Decode(details)
	if err != nil {
		log.Warningf(gen.ctx, "Unable to decode name meaning: %#v", err)
		return details, nil
	}

	//Convert m/f/mf to "Boy", "Girl", "Either"
	for i, detail := range details.Usages {
		details.Usages[i].UsageGender = gen.convertGenderCode(detail.UsageGender)
	}

	gen.addNameToStore(details)

	return details, nil
}

// convertGenderCode converts f/m/mf to Girl, Boy, Either
func (gen *NameGenerator) convertGenderCode(in string) string {
	if in == "f" {
		return "Girl"
	} else if in == "m" {
		return "Boy"
	} else {
		return "Either"
	}
}

func (gen *NameGenerator) addNameToStore(details *NameDetails) error {

	key := datastore.NewIncompleteKey(gen.ctx, NameDetailsEntityType, nil)
	if _, err := datastore.Put(gen.ctx, key, details); err != nil {
		return err
	}
	return nil
}

func (gen *NameGenerator) randomUsage() string {
	usages := []string{"bibl", "bre", "cela", "celm", "cor", "eng", "fre",
		"ger", "iri", "ita", "nor", "roma", "romm", "sca", "sco",
		"sct", "usa", "wel"}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	result := usages[r.Intn(len(usages)-1)]
	return result
}
