package persist

import (
	"encoding/csv"
	"github.com/AndyNortrup/baby-namer/names"
	//"io/ioutil"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"os"
	"path"
	"strconv"
	"sync"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
	"io/ioutil"
	"net/url"
	"strings"
)

var wg sync.WaitGroup

func LoadNames(ctx context.Context, dir string) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Errorf(ctx, "action=load_names %v", err)
	}

	input := make(chan *names.Name)

	go receiveNames(ctx, input)

	for _, file := range files {
		if file.Mode().IsRegular() && strings.Contains(file.Name(), ".txt") {
			statsFile := path.Join(dir, file.Name())
			urlVal := url.Values{}
			urlVal.Add(formValueFilePath, statsFile)
			t := taskqueue.NewPOSTTask(HandleReadStatsTask, urlVal)
			_, err := taskqueue.Add(ctx, t, "loaddata")
			if err != nil {
				log.Errorf(ctx, "action=LoadNames task_name=%v error=%v", statsFile, err)
			} else {
				log.Infof(ctx, "action=LoadNames task_name=%v result=queued", statsFile)
			}
		}

	}

	wg.Wait()
	close(input)
}

func deleteNamesAndStats(ctx context.Context) error {
	//Delete existing names and stats
	err := deleteNameData(ctx)
	if err != nil {
		log.Errorf(ctx, "method=load_names, action=delete_names, error=%v", err)
		return err
	}

	err = deleteStatsData(ctx)
	if err != nil {
		log.Errorf(ctx, "method=load_names, action=delete_stats, error=%v", err)
		return err
	}

	return nil
}

func readStatsFile(ctx context.Context, path string) {
	//defer wg.Done()

	log.Infof(ctx, "action=readStatsFile status=start_import file=%v", path)
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return
	}

	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return
	}

	lines, err := readLines(file)
	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return
	}

	convertLinesToStat(ctx, lines)
	log.Infof(ctx, "action=readStatsFile status=finished_import file=%v", path)
}

func readLines(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	return reader.ReadAll()
}

//Converts an array representing a statistic to a SSANameStat and adds the resulting name to the channel to be merged.
// The line array has three fields:
// [0] = Name
// [1] = Gender (M/F)
// [2] = Number of occurrences that year
// [3] = Rank for that year
// [4] = Year

func convertLinesToStat(ctx context.Context, lines [][]string) {

	wg := &sync.WaitGroup{}

	for _, line := range lines {
		wg.Add(1)

		go func(ctx context.Context, line []string) {
			defer wg.Done()
			name := names.NewName(line[0], names.GetGender(line[1]))
			stat := extractStatFromLine(line)
			dStat := NewDatastoreStat(name, *stat)

			dName := newDatastoreName(name)
			_, err := nds.Put(ctx, datastore.NewKey(ctx, entityTypeName, name.Key(), 0, nil), dName)
			if err != nil {
				log.Errorf(ctx, "action=convertLinesToStat err=%v", err)
			}
			_, err = nds.Put(ctx, dStat.newStatKey(ctx), dStat)
			if err != nil {
				log.Errorf(ctx, "action=convertLinesToStat err=%v", err)
			}
		}(ctx, line)
	}

	wg.Wait()
}

func extractStatFromLine(line []string) *names.Stat {
	year, _ := strconv.Atoi(line[4])
	occurrence := extractOccurrences(line)
	rank := extractRank(line)
	stat := names.NewNameStat(year, rank, occurrence)
	return stat
}

// extractOccurrences pulls the occurrence field from the line and converts it to an integer.  Returns zero in the case
// of a failure.
func extractOccurrences(line []string) int {
	occurrence, err := strconv.Atoi(line[2])
	if err != nil {
		occurrence = 0
	}
	return occurrence
}

func extractRank(line []string) int {
	occurrence, err := strconv.Atoi(line[3])
	if err != nil {
		occurrence = 0
	}
	return occurrence
}

func receiveNames(ctx context.Context, in <-chan *names.Name) {
	//data := NewDatastoreManager(ctx)

	for name := range in {
		dName := newDatastoreName(name)

		//err := data.AddName(name)
		key := datastore.NewKey(ctx, entityTypeName, name.Key(), 0, nil)
		_, err := nds.Put(ctx, key, dName)
		if err != nil {
			log.Errorf(ctx, "Action:load_data: %v", err)
		}
	}
	log.Infof(ctx, "--------------Finished importing names!-----------------")
}
