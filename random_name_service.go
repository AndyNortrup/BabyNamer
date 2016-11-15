package main

import (
	"encoding/xml"
	"math/rand"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

const randomNameEndpoint string = "http://www.behindthename.com/api/random.php?"

type RandomNameService struct {
	ctx    context.Context
	client *http.Client
}

func NewRandomNameService(ctx context.Context) *RandomNameService {
	service := &RandomNameService{}

	service.ctx = ctx
	service.client = urlfetch.Client(ctx)

	return service
}

func (u *RandomNameService) getNameFromService() (*NameDetails, error) {
	addr := randomNameEndpoint +
		"usage=" + u.randomUsage() +
		"&key=an468794&number=1&gender=f"
	log.Infof(u.ctx, "Name request API: %v", addr)
	nameReq, err := u.client.Get(addr)

	if err != nil {
		log.Warningf(u.ctx, "Error getting random name: %v\n", err)
		return nil, err
	}

	defer nameReq.Body.Close()

	decoder := xml.NewDecoder(nameReq.Body)
	name := &BabyName{}
	err = decoder.Decode(name)

	log.Infof(u.ctx, "Recieved '%#v' from service.", name)
	return u.getNameDetails(name)
}

func (gen *RandomNameService) getNameDetails(name *BabyName) (*NameDetails, error) {
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

	mgr := &DatastoreNameManager{
		ctx: gen.ctx,
	}

	mgr.addNameToStore(details)
	log.Infof(gen.ctx, "Name details recieved for '%#v'", details)

	return details, nil
}

func (gen *RandomNameService) randomUsage() string {
	usages := []string{"eng", "iri", "sco", "fre", "wel"}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	return usages[r.Intn(len(usages)-1)]
}

// convertGenderCode converts f/m/mf to Girl, Boy, Either
func (gen *RandomNameService) convertGenderCode(in string) string {
	if in == "f" {
		return "Girl"
	} else if in == "m" {
		return "Boy"
	} else {
		return "Either"
	}
}
