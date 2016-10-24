package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const SettingUsageEntityType string = "SettingUsage"

type SettingUsage struct {
	NameOrigin
	Enabled bool   `datastore:"Enabled"`
	User    string `datastore:"User"`
}

func NewSettingUsage(code string, enabled bool, user string) *SettingUsage {
	origin := getNameOrigins()[code]
	return &SettingUsage{NameOrigin: origin, Enabled: enabled, User: user}
}

func updateUsage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	code := mux.Vars(r)["code"]
	enabled, err := strconv.ParseBool(mux.Vars(r)["status"])

	if err != nil {
		log.Errorf(ctx, "Failed to convert status to boolean: %v", err)
	}

	user := user.Current(ctx)

	setting := NewSettingUsage(code, enabled, user.Email)

	log.Infof(ctx, "Recieved request to change usage=%v to %v for %v",
		setting.Code, setting.Enabled, setting.User)

	mgr := NewSettingUsageDatastoreManager(ctx)
	mgr.setUsage(setting)
}

func getAllUserUsages(user string, ctx context.Context) map[string]*SettingUsage {
	result := make(map[string]*SettingUsage)
	query := datastore.NewQuery(SettingUsageEntityType).
		Filter("User =", user).Filter("Enabled =", true)

	setting := &SettingUsage{}
	for t := query.Run(ctx); ; {
		_, err := t.Next(setting)

		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "Failed to get settings from datastore: %v", err)
		} else {
			result[setting.Code] = setting
		}
	}
	log.Infof(ctx, "Enabled user settings for %v: %v", user, result)
	return result
}
