package persist

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

var formValueFilePath string = "filePath"
var HandleReadStatsTask string = "/admin/read_stats_file"

func HandleLoadData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)

	LoadNames(ctx, "names")
	w.WriteHeader(http.StatusOK)
}

func HandleReadStatsFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Starting read stats file: %v", r.FormValue(formValueFilePath))
	readStatsFile(ctx, r.FormValue(formValueFilePath))
}
