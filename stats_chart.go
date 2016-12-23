package babynamer

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/persistance"
	"github.com/gorilla/mux"
	"github.com/wcharczuk/go-chart"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

type StatsChart struct {
	Name *names.Name
}

func NewStatsChart(name *names.Name) *StatsChart {
	return &StatsChart{Name: name}
}

func (sc *StatsChart) RenderChart(w http.ResponseWriter) {

	//Get the stats and put them in arrays to use as time series.
	dates := []time.Time{}
	occurrences := []float64{}
	rank := []float64{}
	for _, stat := range sc.Name.SortedStats() {
		dates = append(dates, stat.YearAsTime())
		occurrences = append(occurrences, float64(stat.Occurrences))
		rank = append(rank, float64(stat.Rank))
	}

	graph := &chart.Chart{
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxisSecondary: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Width:  300,
		Height: 300,

		Series: []chart.Series{
			chart.TimeSeries{
				XValues: dates,
				YValues: occurrences,
				Name:    "Occurrences",
				YAxis:   chart.YAxisPrimary,
			},
			chart.TimeSeries{
				XValues: dates,
				YValues: rank,
				Name:    "Rank",
				YAxis:   chart.YAxisSecondary,
			},
		},
	}

	w.Header().Set("Content-Type", "images/png")
	graph.Render(chart.PNG, w)
}

func handleChart(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	nameStr := mux.Vars(r)["name"]
	genderStr := mux.Vars(r)["gender"]
	gender := names.GetGender(genderStr)
	log.Infof(ctx, "Action=StatsChart Name=%v Gender=%s", nameStr, gender.GoString())

	//get name from Datastore
	data := persist.NewDatastoreManager(ctx)
	name, err := data.GetName(nameStr, gender)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c := NewStatsChart(name)
	c.RenderChart(w)
}
