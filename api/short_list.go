package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type ShortList struct {
	Names []*NameDetails
}

func NewShortList(ctx context.Context) ShortList {
	list := ShortList{
		Names: []*NameDetails{},
	}
	list.getList(ctx)
	return list
}

func (list *ShortList) getList(ctx context.Context) error {
	query := datastore.NewQuery(EntityTypeNameDetails).
		Filter("ApprovedBy >", "").
		Filter("RejectedBy =", "")

	for i := query.Run(ctx); ; {
		detail := &NameDetails{}
		_, err := i.Next(detail)
		if err == datastore.Done {
			return nil
		} else if err != nil {
			log.Errorf(ctx, "Error getting items from short list: %v", err)
			return err
		} else {
			log.Infof(ctx, "Adding name: %v to short list", detail)
			list.Names = append(list.Names, detail)
		}
	}
}
