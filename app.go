package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const EntityTypeNameDetails string = "NameDetails"

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
	r.HandleFunc("/short-list", handleShortList)
	r.HandleFunc("/settings", handleSettings)
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("templates/fonts"))))

	http.Handle("/", r)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
