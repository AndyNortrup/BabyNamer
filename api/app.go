package main

import (
	"net/http"

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
	// api := r.PathPrefix("/api").Subrouter()
	// // api/random-name
	// api.HandleFunc("/random-name", getNameHandler)
	// api.HandleFunc("/get-supported-usages", getAllUsages)
	// api.HandleFunc("/make-decision/{name}/{choice:(yes|no)}", apiMakeDecision)
	// api.HandleFunc("/update-usage/{code}/{status:(true|false)}", updateUsage)

	r.HandleFunc("/", namesPage)
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("templates/fonts"))))

	http.Handle("/", r)
}

// func getNameHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	ctx := appengine.NewContext(r)
// 	gen := NewNameGenerator(ctx)
// 	gen.getName(w, r)
// }
//
// func apiMakeDecision(w http.ResponseWriter, r *http.Request) {
//
// 	defer r.Body.Close()
//
// 	//Get the user's name
// 	ctx := appengine.NewContext(r)
// 	username := user.Current(ctx)
//
// 	name := mux.Vars(r)["name"]
// 	var decision bool
//
// 	if mux.Vars(r)["choice"] == "yes" {
// 		decision = true
// 	} else {
// 		decision = false
// 	}
//
// 	log.Infof(ctx, "User: %v says %v to %v", username, decision, name)
// 	err := updateNameRecommendations(name, username.Email, decision, ctx)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
//
// 	log.Infof(ctx, "Requesting new name to return")
// 	http.Redirect(w, r, "/?previous="+name, http.StatusSeeOther)
// }
