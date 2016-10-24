package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"github.com/gorilla/mux"
)

const NameDetailsEntityType string = "NameDetails"

type BabyName struct {
	Name   string `xml:"names>name"`
	Gender string
}

type Usage struct {
	UsageFull   string `xml:"usage>usage_full"`
	UsageGender string `xml:"usage>usage_gender"`
}

func init() {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api").Subrouter()
	// api/random-name
	sr.HandleFunc("/random-name", getNameHandler)
	sr.HandleFunc("/get-supported-usages", getAllUsages)
	sr.HandleFunc("/make-decision/{name}/{choice:(yes|no)}", makeDecision)
	sr.HandleFunc("/update-usage/{code}/{status:(true|false)}", updateUsage)
	http.Handle("/", r)
}

func getNameHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	getName(w, r)
}

func makeDecision(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	//Get the user's name
	ctx := appengine.NewContext(r)
	username := user.Current(ctx)
	name := mux.Vars(r)["name"]
	var decision bool

	if mux.Vars(r)["choice"] == "yes" {
		decision = true
	} else {
		decision = false
	}

	log.Infof(ctx, "User: %v says %v to %v", username, decision, name)

	//Get the name from the datastore
	query := datastore.NewQuery(NameDetailsEntityType).
		Filter("Name =", name)

	for t := query.Run(ctx); ; {
		details := &NameDetails{}
		key, err := t.Next(details)

		if err == datastore.Done {
			log.Infof(ctx, "Couldn't find name in datastore.")
			break
		}
		if err != nil {
			log.Warningf(ctx, "Error retriving Name Details %#v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			err := recordDecision(key, decision, details, username, ctx)
			if err != nil {
				log.Warningf(ctx, "Error updating name with decision.")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			log.Infof(ctx, "Name record updated.")
			break
		}
	}
	log.Infof(ctx, "Requesting new name to return")
	getName(w, r)
}

func getName(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	gen := NewNameGenerator(ctx)
	details, err := gen.getName(mux.Vars(r)["name"])

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Debugf(ctx, "Error getting name: %v", err)
	}

	log.Debugf(ctx, "Sending name: %v", details)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(details)
}

func recordDecision(key *datastore.Key,
	decision bool,
	details *NameDetails,
	username *user.User,
	ctx context.Context) error {

	if decision {
		//Check if this has already been updated
		if details.RecommendedBy == "" {
			details.RecommendedBy = username.String()
		} else {
			details.ApprovedBy = username.String()
		}
	} else {
		details.RejectedBy = username.String()
	}

	_, err := datastore.Put(ctx, key, details)
	return err
}
