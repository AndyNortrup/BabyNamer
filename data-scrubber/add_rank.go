package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		log.Printf("Use: data-scrubber dir")
		os.Exit(1)
	}

	dir := os.Args[1]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Unable to read files in directory. Error=%v", err.Error())
		os.Exit(3)
	}

	for _, file := range files {
		log.Printf("Reading file=%v", file.Name())
		fileDesc, err := os.Open(dir + "/" + file.Name())

		defer fileDesc.Close()
		if err != nil {
			log.Printf("Unable to read file=%v Error=%v", file.Name(), err.Error())
			os.Exit(4)
		}
		lines := csv.NewReader(fileDesc)
		lines.FieldsPerRecord = 3
		linesArray, err := lines.ReadAll()
		if err != nil {
			log.Printf("Failed to read csv=%v error=%v", file.Name(), err.Error())
		}

		out, err := os.Create(os.Args[2] + "/" + file.Name())
		defer out.Close()

		outFile := csv.NewWriter(out)
		m, f := 0, 0
		for _, line := range linesArray {
			occurances, _ := strconv.Atoi(line[2])
			if occurances > 100 {
				if line[1] == "M" {
					m++
					line = append(line, strconv.Itoa(m))
				} else {
					f++
					line = append(line, strconv.Itoa(f))
				}
				outFile.Write(line)
			}
		}
		outFile.Flush()

	}
}
