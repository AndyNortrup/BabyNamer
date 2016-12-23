package persist

import (
	"google.golang.org/appengine"
	"net/http"
)

func HandleLoadData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)

	LoadNames(ctx)
	w.WriteHeader(http.StatusOK)
}
