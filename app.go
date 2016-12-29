package babynamer

import (
	"github.com/AndyNortrup/baby-namer/persistance"
	"github.com/gorilla/mux"
	"net/http"
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
	r.HandleFunc("/settings", handleSettingsPage)
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("templates/fonts"))))

	r.HandleFunc("/admin/load_ssa_data", persist.HandleLoadData)
	r.HandleFunc("/charts/{name}/{gender}", handleChart)

	http.Handle("/", r)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
