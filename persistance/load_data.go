package persist

import (
	"encoding/csv"
	"github.com/AndyNortrup/baby-namer/names"
	"io/ioutil"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func LoadNames(ctx context.Context, dir string) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Errorf(ctx, "action=load_names %v", err)
	}

	input := make(chan *names.Name)

	//Try 4 because we have four cores
	go receiveNames(ctx, input)
	go receiveNames(ctx, input)
	go receiveNames(ctx, input)
	go receiveNames(ctx, input)

	for _, file := range files {
		if file.Mode().IsRegular() && strings.Contains(file.Name(), ".txt") {
			statsFile := path.Join(dir, file.Name())
			wg.Add(1)
			go readStatsFile(ctx, statsFile, input)
		}

	}

	wg.Wait()
	close(input)
}

func readStatsFile(ctx context.Context, path string, ch chan<- *names.Name) {
	defer wg.Done()

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return
	}

	year, err := convertFileNameToYear(file.Name())
	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return
	}

	lines, err := readLines(file)
	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		return
	}

	convertLinesToStat(lines, year, ch)
	log.Infof(ctx, "action=readStatsFile status=finished_import file=%v", path)
}

func convertFileNameToYear(p string) (int, error) {
	_, file := path.Split(p)
	return strconv.Atoi(file[3:7])
}

func readLines(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	return reader.ReadAll()
}

//Converts an array representing a statistic to a SSANameStat and adds the resulting name to the channel to be merged.
// The line array has three fields:
// [0] = Name
// [1] = Gender (M/F)
// [2] = Number of occurrences that year

func convertLinesToStat(lines [][]string, year int, out chan<- *names.Name) {
	var m, f int
	for _, line := range lines {
		name := names.NewName(line[0], names.GetGender(line[1]))
		occurrence := extractOccurrences(line)
		if line[1] == "M" {
			m++
			name.AddStat(names.NewNameStat(year, m, occurrence))
		} else {
			f++
			name.AddStat(names.NewNameStat(year, f, occurrence))
		}
		out <- name
	}
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

func receiveNames(ctx context.Context, in <-chan *names.Name) {
	data := NewDatastoreManager(ctx)

	for name := range in {
		err := data.AddName(name)
		if err != nil {
			log.Errorf(ctx, "Action:load_data: %v", err)
		}
	}
	log.Infof(ctx, "--------------Finished importing names!-----------------")
}
