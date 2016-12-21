package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"net/http"
)

func HandleLoadData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)
	mgr := NewDatastoreManager(ctx)

	names := LoadNames()

	err := addNamesToDatastore(ctx, names, mgr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func addNamesToDatastore(ctx context.Context, names map[string]*names.Name, mgr DataManager) error {
	for _, value := range names {
		mgr.AddName(value)
	}
	return nil
}
