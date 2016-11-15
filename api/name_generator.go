package main

import (
	"encoding/xml"
	"errors"
	"math/rand"
	"time"
	//"math/rand"
	"net/http"
	//"time"

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

const randomNameEndpoint string = "http://www.behindthename.com/api/random.php?"

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

	if name != nil && name.Name != "" {
		log.Infof(gen.ctx, "Returning recommended name: %v", name.Name)
		return name, nil
	}

	log.Infof(gen.ctx, "Returning random name")
	isRejected := true
	var randomName NameDetails
	for isRejected {

		randomName, err := gen.getNameFromService()
		if err != nil || randomName.Name == "" {
			log.Warningf(gen.ctx, "Error retriving random name: %v", err)
		} else if !gen.isRejected(randomName) {
			return randomName, nil
		}
	}

	return &randomName, nil
}

//getRecommendedName gets a list of names from the datastore that were
// recommended by a different user, and have not been rejected by any user
func (gen *NameGenerator) getRecommendedName(previous string) (*NameDetails, error) {
	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("ApprovedBy =", "").
		Filter("RecommendedBy <", "").
		Filter("RejectedBy =", "")

	url, urlErr := user.LogoutURL(gen.ctx, "/")
	if urlErr != nil {
		log.Errorf(gen.ctx, "Error building logout URL: %v", url)
	}

	for t := query.Run(gen.ctx); ; {
		details := &NameDetails{}
		_, err := t.Next(details)
		if err == datastore.Done {
			log.Infof(gen.ctx, "No recommended names found.")
			return nil, nil
		} else if err != nil {
			log.Infof(gen.ctx, "Error retriving recommended names: %v", err)
			return nil, err
		} else if details.RecommendedBy != gen.user && previous != details.Name {
			log.Infof(gen.ctx, "Returning recommended name: %v", details)
			return details, nil
		}
	}
}

func (gen *NameGenerator) getNameFromService() (*NameDetails, error) {
	addr := randomNameEndpoint +
		"usage=" + gen.randomUsage() +
		"&key=an468794&number=1&gender=f"
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

	log.Infof(gen.ctx, "Recieved '%#v' from service.", name)
	return gen.getNameDetails(name)
}

func (gen *NameGenerator) getNameDetails(name *BabyName) (*NameDetails, error) {
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
	log.Infof(gen.ctx, "Name details recieved for '%#v'", details)

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

	if details.Name == "" {
		return errors.New("Can't write blank name to datastore.")
	}

	//Check if the name already exists by searching for it.
	key, _ := gen.getKeyForName(details)
	if key != nil {
		return errors.New("Name already exists in datastore.")
	}

	key = datastore.NewIncompleteKey(gen.ctx, EntityTypeNameDetails, nil)
	if _, err := datastore.Put(gen.ctx, key, details); err != nil {
		log.Warningf(gen.ctx, "Error writing name to datastore: %v", err)
		return err
	}
	return nil
}

func (gen *NameGenerator) randomUsage() string {
	usages := []string{"eng", "iri", "sco", "fre", "wel"}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	return usages[r.Intn(len(usages)-1)]
}

//getKeyForName returns the datastore key for a given name value. Use this to
// get a key so you update a name's information rather than write a new copy.
func (gen *NameGenerator) getKeyForName(details *NameDetails) (*datastore.Key, error) {
	t := datastore.NewQuery(EntityTypeNameDetails).
		Filter("Name =", details.Name).
		Run(gen.ctx)

	for {
		results := &NameDetails{}
		key, err := t.Next(results)

		if err == datastore.Done {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}
		return key, nil
	}
}

//isRejected returns a boolean value to tell you if a name generated by the
// random name service has been previously rejected
func (gen *NameGenerator) isRejected(details *NameDetails) bool {

	//Get all of the instnace of the name that have been rejected by someone.
	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("Name =", details.Name).
		Filter("RejectedBy >", "")

	t := query.Run(gen.ctx)

	for {
		result := &NameDetails{}
		_, err := t.Next(result)
		if err == datastore.Done {
			//If we get no result than the name has not been rejected.
			return false
		}
		if err != nil {
			//If we get an error, we'll just say no it hasn't been rejected
			log.Warningf(gen.ctx,
				"Unable to validate if name has already been rejected: %v", err)
			return false
		}
		//If there are any results at all, than the name has been rejected
		log.Infof(gen.ctx, "The name '%v' has not been previously rejected.",
			details.Name)
		return true
	}
}
