package babynamer

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/settings"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"sort"
)

func getShortList(ctx context.Context) (names.NameList, error) {

	oldList, err := getDeprecatedShortList(ctx)
	if err != nil {
		return nil, err
	}

	shortList, err := getPersistenceShortList(ctx)
	if err != nil {
		return nil, err
	}

	shortList = append(shortList, oldList...)
	sort.Sort(shortList)

	return shortList, nil

}

func getPersistenceShortList(ctx context.Context) (names.NameList, error) {
	mgr := newDataManager(ctx)
	usr := user.Current(ctx)
	partner, err := settings.GetPartner(ctx, usr)
	if err != nil {
		log.Errorf(ctx, "action=getShortList step=getPartner err=%v", err)
		return nil, err
	}
	shortList, err := mgr.GetShortList(usr, &user.User{Email: partner.PartnerEmail})
	if err != nil {
		log.Errorf(ctx, "action=getShortList step=GetShortList err=%v", err)
		return nil, err
	}
	return shortList, nil
}

func getDeprecatedShortList(ctx context.Context) (names.NameList, error) {
	list := []*NameDetails{}
	_, err := datastore.NewQuery(EntityTypeNameDetails).
		Filter("ApprovedBy >", "").
		Filter("RejectedBy =", "").GetAll(ctx, &list)

	outList := names.NameList{}
	for _, value := range list {
		outList = append(outList, &names.Name{
			Name:   value.Name,
			Gender: names.GetGender(value.Gender),
		})
	}

	return outList, err
}
