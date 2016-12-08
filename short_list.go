package babynamer

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"sort"
)

type ShortList []*NameDetails

func getShortList(ctx context.Context) (ShortList, error) {
	list := make(ShortList, 0)

	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("ApprovedBy >", "").
		Filter("RejectedBy =", "")

	for i := query.Run(ctx); ; {
		detail := &NameDetails{}
		_, err := i.Next(detail)
		if err == datastore.Done {
			sort.Sort(list)
			return list, nil
		} else if err != nil {
			return nil, err
		} else {
			list = append(list, detail)
		}
	}
}

func (list ShortList) Len() int {
	return len(list)
}

func (list ShortList) Less(i, j int) bool {
	return list[i].Name < list[j].Name
}

func (list ShortList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
