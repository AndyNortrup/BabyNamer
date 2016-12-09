package babynamer

import (
	"net/http"

	"github.com/gorilla/mux"
)

const EntityTypeNameDetails string = "NameDetails"

type BabyName struct {
	Name   string `xml:"names>name"`
	Gender string
}

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/", namesPage)
	r.HandleFunc("/short-list", handleShortList)
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("templates/fonts"))))

	http.Handle("/", r)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
