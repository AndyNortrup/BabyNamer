package settings

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const entityTypePartner string = "partner"
const filterSelf string = "SelfEmail ="

type Partner struct {
	SelfEmail    string
	PartnerEmail string
}

func GetPartner(ctx context.Context, self *user.User) (*Partner, error) {
	query := buildGetSettingsQuery(self)
	partner := []*Partner{}
	_, err := query.GetAll(ctx, &partner)
	if err != nil {
		log.Errorf(ctx, "action=GetPartner error=%v", err)
		return nil, err
	}
	if len(partner) > 0 {
		return partner[0], nil
	}
	return nil, nil
}

func buildGetSettingsQuery(self *user.User) *datastore.Query {
	return datastore.NewQuery(entityTypePartner).Filter(filterSelf, self.Email)
}

func SetPartner(ctx context.Context, self *user.User, partnerEmail string) error {

	log.Infof(ctx, "action=SetPartner event=starting set partner")
	query := buildGetSettingsQuery(self).KeysOnly()

	partner := []*Partner{}
	keys, err := query.GetAll(ctx, &partner)
	if err != nil {
		log.Errorf(ctx, "action=GetPartner error=%v", err)
		return err
	}

	var writeKey *datastore.Key
	if len(keys) > 0 {
		writeKey = keys[0]
	} else {
		writeKey = datastore.NewIncompleteKey(ctx, entityTypePartner, nil)
	}

	output := &Partner{SelfEmail: self.Email, PartnerEmail: partnerEmail}
	_, err = datastore.Put(ctx, writeKey, output)
	if err != nil {
		log.Errorf(ctx, "action=GetPartner error=%v", err)
		return err
	}
	log.Infof(ctx, "action=SetPartner successfully set partner value=%v", output)
	return nil
}
