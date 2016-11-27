package main

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type SettingUsageDatastoreManager struct {
	ctx context.Context
}

func NewSettingUsageDatastoreManager(ctx context.Context) *SettingUsageDatastoreManager {
	return &SettingUsageDatastoreManager{ctx: ctx}
}

func (mgr *SettingUsageDatastoreManager) setUsage(setting *SettingUsage) {
	//get key that matches existing setting for that usage
	key := mgr.makeKey(setting)

	//write new value to that key
	_, err := datastore.Put(mgr.ctx, key, setting)
	if err != nil {
		log.Errorf(mgr.ctx, "Failed to write setting to datastore: %v", err)
	}
	log.Infof(mgr.ctx, "Updated setting key: %v", key)
}

//getSettingsKey searches the datastore for existing settings for this usage and
// user.  If none is found it returns a new key to be written to, if the setting
// exists then the existing key is returned.
func (mgr *SettingUsageDatastoreManager) getSettingsKey(setting *SettingUsage) *datastore.Key {
	query := datastore.NewQuery(SettingUsageEntityType).
		Filter("User =", setting.User).
		Filter("Code =", setting.Code).
		KeysOnly()

	for t := query.Run(mgr.ctx); ; {
		key, err := t.Next(&SettingUsage{})

		//There is no key, so this is the first time we've set this value
		if err == datastore.Done {
			newKey := mgr.makeKey(setting)
			log.Infof(mgr.ctx, "Returning new Usage Setting Key: %v", newKey)
			return newKey
		}
		if err != nil {
			log.Errorf(mgr.ctx, "Failed to get settings keys: %v", err)
			return mgr.makeKey(setting)
		}
		//Get out of the loop we have our key
		log.Infof(mgr.ctx, "Returning existing Usage Setting Key: %v", key)
		return key
	}
}

func (mgr *SettingUsageDatastoreManager) makeKey(setting *SettingUsage) *datastore.Key {
	return datastore.NewKey(mgr.ctx,
		SettingUsageEntityType,
		setting.User+"-"+setting.Code,
		0, nil)
}
