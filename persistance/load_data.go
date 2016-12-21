package persist

import (
	"encoding/csv"
	"github.com/AndyNortrup/baby-namer/names"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func LoadNames() map[string]*names.Name {

	dir := "names"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	input := make(chan *names.Name)
	result := make(chan map[string](*names.Name))

	go receiveNames(input, result)

	for _, file := range files {

		if file.Mode().IsRegular() {
			statsFile := path.Join(dir, file.Name())
			wg.Add(1)
			go readStatsFile(statsFile, input)
		}
	}
	wg.Wait()
	close(input)

	end := <-result
	return end
}

func readStatsFile(path string, ch chan<- *names.Name) {
	defer wg.Done()

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	year, err := convertFileNameToYear(file.Name())
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	lines, err := readLines(file)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	convertLinesToStat(lines, year, ch)
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
		name := names.NewName(line[0], line[1])
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

func receiveNames(in <-chan *names.Name, out chan<- map[string]*names.Name) {
	result := make(map[string]*names.Name)

	for name := range in {
		result = addNameToMasterList(result, name)
	}
	out <- result
	close(out)
}

func addNameToMasterList(names map[string]*names.Name, name *names.Name) map[string]*names.Name {

	key := name.Key()

	if names[key] == nil {
		names[key] = name
	} else {
		update := names[key]
		for _, stat := range name.Stats {
			update.AddStat(stat)
			names[key] = update
		}
	}

	return names
}
