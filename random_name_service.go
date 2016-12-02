package main

import (
	"encoding/xml"

	"golang.org/x/net/context"

	"net/http"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const randomNameEndpoint string = "http://www.behindthename.com/api/random.php?"
const nameDetailsEndpoint string = "http://www.behindthename.com/api/lookup.php?"

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

func (u *RandomNameService) getNameFromService(usage string) (*NameDetails, error) {

	addr := randomNameEndpoint +
		"usage=" + usage +
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
	if err != nil {
		log.Errorf(u.ctx, "Failed to decode name from service: %v", err)
	}

	log.Infof(u.ctx, "Recieved '%#v' from service.", name)
	details, err := u.getNameDetails(name)
	if err != nil {
		log.Errorf(u.ctx, "Error retriving name details: %v", err)
	}
	return details, err
}

func (gen *RandomNameService) getNameDetails(name *BabyName) (*NameDetails, error) {
	details := &NameDetails{}
	address := nameDetailsEndpoint + "key=an468794&name=" + name.Name
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

	log.Infof(gen.ctx, "Name details recieved for '%#v'", details)
	err = mgr.addNameToStore(details)

	if err == DuplicateNameError {
		log.Infof(mgr.ctx, err.Error())
	} else if err != nil {
		log.Errorf(mgr.ctx, "Failed to write name to datastore: %v", err)
		return nil, err
	}

	return details, nil
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
