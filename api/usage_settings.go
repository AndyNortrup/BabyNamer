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
	origin := NameOrigin{Code: code}
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
	setting.setUsage(ctx)
}

func (setting SettingUsage) setUsage(ctx context.Context) {
	//get key that matches existing setting for that usage
	key := setting.getSettingsKey(ctx)

	//write new value to that key
	datastore.Put(ctx, key, setting)
	log.Infof(ctx, "Updated setting key: %v", key)
}

func (setting SettingUsage) getSettingsKey(ctx context.Context) *datastore.Key {
	query := datastore.NewQuery(SettingUsageEntityType)
	query.Filter("User =", setting.User).
		Filter("Code =", setting.Code).
		KeysOnly()

	for t := query.Run(ctx); ; {
		key, err := t.Next(&SettingUsage{})

		//There is no key, so this is the first time we've set this value
		if err == datastore.Done {
			return datastore.NewIncompleteKey(ctx, SettingUsageEntityType, nil)
		}
		if err != nil {
			log.Errorf(ctx, "Failed to get settings keys: %v", err)
			return datastore.NewIncompleteKey(ctx, SettingUsageEntityType, nil)
		}
		//Get out of the loop we have our key
		return key
	}
}

func getAllUserUsages(user string, ctx context.Context) map[string]*SettingUsage {
	result := make(map[string]*SettingUsage)
	query := datastore.NewQuery(SettingUsageEntityType).
		Filter("User =", user)

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
	return result
}
