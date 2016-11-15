package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func handleSettings(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)
	user := user.Current(ctx)

	settingManager := NewSettingUsageDatastoreManager(ctx)
	settingManager.setUsage(&SettingUsage{
		NameOrigin: NameOrigin{Code: "USA", Plain: "American"},
		Enabled:    true,
		User:       "test",
	})

	usages := NewUsageList(ctx, user.String())
	log.Debugf(ctx, "Usage List for user: %v", user.String())

	t, err := getSettingTemplate()
	if err != nil {
		log.Errorf(ctx, "Failed to load settings template: %v", err)
	}
	t.Execute(w, usages)
}
