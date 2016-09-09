package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"

	"github.com/gorilla/mux"
)

const NameDetailsEntityType string = "NameDetails"

func init() {
	r := mux.NewRouter()

	sr := r.PathPrefix("/api").Subrouter()
	// api/random-name
	sr.HandleFunc("/random-name", getName)
	sr.HandleFunc("/sendToMaybe/{name}", sendToMaybe)

	r.HandleFunc("/{rest:.*}", serveStatic)
	http.Handle("/", r)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	log.Printf("path: %v", r.URL.Path)
	http.ServeFile(w, r, "static/"+r.URL.Path)
}

func getName(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	nameReq, err := client.Get(
		"http://www.behindthename.com/api/random.php?key=an468794&number=1")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error getting random name: %v\n", err)
	}

	defer nameReq.Body.Close()

	decoder := xml.NewDecoder(nameReq.Body)
	names := &BabyName{}
	err = decoder.Decode(names)

	details, err := getNameInfo(names, client)
	if err != nil {
		http.Error(w, "Internal server Error", 500)
	}

	err = addNameToStore(details, ctx)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(details)
}

func getNameInfo(name *BabyName, client *http.Client) (*NameDetails, error) {
	details := &NameDetails{}
	address := "http://www.behindthename.com/api/lookup.php?key=an468794&name=" +
		name.Name
	nameDetailsReq, err :=
		client.Get(address)

	if err != nil {
		log.Printf("Unable to retrieve name details: %#v", err)
		return details, err
	}

	defer nameDetailsReq.Body.Close()
	decoder := xml.NewDecoder(nameDetailsReq.Body)

	err = decoder.Decode(details)
	if err != nil {
		log.Printf("Unable to decode name meaning: %#v", err)
		return details, nil
	}

	//Convert m/f/mf to "Boy", "Girl", "Either"
	for i, detail := range details.Usages {
		if detail.UsageGender == "f" {
			details.Usages[i].UsageGender = "Girl"
		} else if detail.UsageGender == "m" {
			details.Usages[i].UsageGender = "Boy"
		} else {
			details.Usages[i].UsageGender = "Either"
		}
	}

	return details, nil
}

type BabyName struct {
	// <response>
	// 	<names>
	// 		<name>Jeane</name>
	// 	</names>
	// </response>

	Name   string `xml:"names>name"`
	Gender string
}

type NameDetails struct {
	// <response>
	// 	<name_detail>
	// 		<name>Andy</name>
	//		<gender>mf</gender>
	// 		<usages>
	// 			<usage>
	// 				<usage_code>eng</usage_code>
	// 				<usage_full>English</usage_full>
	// 				<usage_gender>m</usage_gender>
	// 			</usage>
	// 			<usage>
	// 				<usage_code>eng</usage_code>
	// 				<usage_full>English</usage_full>
	// 				<usage_gender>f</usage_gender>
	// 			</usage>
	// 		</usages>
	// 	</name_detail>
	// </response>
	Name          string  `xml:"name_detail>name"`
	Gender        string  `xml:"name_detail>gender"`
	Usages        []Usage `xml:"name_detail>usages"`
	RecommendedBy string
	ApprovedBy    string
	RejectedBy    string
}

type Usage struct {
	UsageFull   string `xml:"usage>usage_full"`
	UsageGender string `xml:"usage>usage_gender"`
}

func addNameToStore(details *NameDetails, context context.Context) error {
	key := datastore.NewIncompleteKey(context, NameDetailsEntityType, nil)
	if _, err := datastore.Put(context, key, details); err != nil {
		return err
	}

	return nil
}

func sendToMaybe(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	//Get the user's name
	log.Println("sendToMaybe called")
	ctx := appengine.NewContext(r)
	username := user.Current(ctx)
	name := mux.Vars(r)["name"]
	log.Printf("User: %v wants to save %v", username, name)

	//Get the name from the datastore
	query := datastore.NewQuery(NameDetailsEntityType).
		Filter("Name =", name)
	for t := query.Run(ctx); ; {
		details := &NameDetails{}
		key, err := t.Next(details)

		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Printf("Error retriving Name Details %#v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			//Check if this has already been updated
			if details.RecommendedBy != "" {
				details.RecommendedBy = username.String()
			} else {
				details.ApprovedBy = username.String()
			}

			_, err = datastore.Put(ctx, key, details)

			if err != nil {
				log.Printf("Unable to update datastore with recommendation: %#v", err)
				http.Error(w, "Internal server Error", 500)
			}
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}
}
