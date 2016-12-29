package settings

import "google.golang.org/appengine/user"

type Setting interface {
	getKey() string
	GetValue() string
}

type User struct {
	user.User
	settings map[string]*Setting `datastore:"-"`
	Settings []*Setting
}
