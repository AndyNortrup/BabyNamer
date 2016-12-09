package babynamer

import (
	"github.com/AndyNortrup/baby-namer/usage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

func scrapeRandomNames(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := appengine.NewContext(r)
	rndService := NewRandomNameService(ctx)
	usageGen := usage.NewUsageGenerator(ctx)

	name, err := rndService.getNameFromService(usageGen.RandomUsageCode())

	if err != nil || name.Name == "" {
		log.Errorf(ctx, "Failed to scrape random name: %v", err)
	}

	dataMgr := NewDatastoreNameManager(ctx, "")
	err = dataMgr.addNameToStore(name)
	if err != nil {
		log.Errorf(ctx, "Failed to write name to datastore: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}
