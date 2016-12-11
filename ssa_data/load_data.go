package ssa_data

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func LoadNames() map[string]*Name {

	dir := "names"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	names := make(chan *Name)
	result := make(chan map[string]*Name)

	go receiveNames(names, result)

	for _, file := range files {

		if file.Mode().IsRegular() {
			statsFile := path.Join(dir, file.Name())
			wg.Add(1)
			go readStatsFile(statsFile, names)
		}
	}
	wg.Wait()
	close(names)

	end := <-result
	return end
}

func readStatsFile(path string, ch chan<- *Name) {
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
// [2] = Number of occurnaces that year

func convertLinesToStat(lines [][]string, year int, out chan<- *Name) {
	var m, f int
	for _, line := range lines {
		name := NewName(line[0], line[1])
		occurrence := extractOccurrences(line)
		if line[1] == "M" {
			m++
			name.addStat(NewSSANameStat(year, m, occurrence))
		} else {
			f++
			name.addStat(NewSSANameStat(year, f, occurrence))
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

func receiveNames(in <-chan *Name, out chan<- map[string]*Name) {
	result := make(map[string]*Name)

	for name := range in {
		result = addNameToMasterList(result, name)
	}
	out <- result
	close(out)
}

func addNameToMasterList(names map[string]*Name, name *Name) map[string]*Name {

	if names[name.Name] == nil {
		names[name.Name] = name
	} else {
		update := names[name.Name]
		for _, stat := range name.Stats {
			update.addStat(stat)
			names[name.Name] = update
		}
	}
	return names
}
