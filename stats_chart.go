package babynamer

import (
	"github.com/AndyNortrup/baby-namer/names"
	"github.com/AndyNortrup/baby-namer/persistance"
	"github.com/gorilla/mux"
	"github.com/wcharczuk/go-chart"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
	"time"
)

type StatsChart struct {
	Name *names.Name
}

func NewStatsChart(name *names.Name) *StatsChart {
	return &StatsChart{Name: name}
}

func (sc *StatsChart) RenderChart(w http.ResponseWriter, ctx context.Context) {

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
		Background: chart.Style{
			Padding: chart.Box{
				Top:   50,
				Right: 50,
			},
		},
		XAxis: chart.XAxis{
			Name:      "Year",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return time.Unix(0, int64(v.(float64))).Format("2006")
			},
		},
		//YAxisSecondary: chart.YAxis {
		//	Name:      "Occurannces",
		//	NameStyle: chart.StyleShow(),
		//	Style: chart.Style{
		//		Show: true,
		//	},
		//	ValueFormatter: func(v interface{}) string {
		//		return strconv.Itoa(int(v.(float64)))
		//	},
		//},
		YAxis: chart.YAxis{
			Name:      "Rank",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return strconv.Itoa(int(v.(float64)))
			},
			Range: &chart.ContinuousRange{
				Descending: true,
				//Min: 1,
				//Max: float64(sc.Name.LowestRank().Year),
			},
		},
		Series: []chart.Series{
			//chart.TimeSeries{
			//	XValues: dates,
			//	YValues: occurrences,
			//	Name:    "Occurrences",
			//	YAxis:   chart.YAxisSecondary,
			//},
			chart.TimeSeries{
				XValues: dates,
				YValues: rank,
				Name:    "Rank",
				YAxis:   chart.YAxisPrimary,
			},
		},
		Width: 1024,
	}

	graph.Elements = []chart.Renderable{
		chart.LegendThin(graph),
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
	c.RenderChart(w, ctx)
}
