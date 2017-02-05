package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var chunkSize = 100

type yearCount struct {
	mRank, fRank, mOccur, fOccur int
}

func main() {
	if len(os.Args) != 3 {
		log.Printf("Use: data-scrubber inDir")
		os.Exit(1)
	}

	inDir := os.Args[1]

	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Printf("Unable to read files in directory. Error=%v", err.Error())
		os.Exit(3)
	}

	processFolder(files, inDir, os.Args[2], chunkSize)
}

func processFolder(files []os.FileInfo, inDir, outDir string, chunkSize int) {
	for _, file := range files {
		log.Printf("Reading file=%v", file.Name())
		fileDesc, err := os.Open(filepath.Join(inDir, file.Name()))

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

		count := &yearCount{}

		for x := 0; x < len(linesArray); x++ {
			outFilePath := filepath.Join(outDir, outFileName(file, x))
			out, err := os.Create(outFilePath)
			if err != nil {
				out.Close()
				continue
			}
			outFile := csv.NewWriter(out)
			limit := setLoopLimit(x, chunkSize, linesArray)
			hasOutput := false

			for ; x < limit; x++ {
				line, use := manipulateLine(linesArray[x], file.Name(), count)
				if use {
					outFile.Write(line)
					hasOutput = true
				}
			}
			outFile.Flush()
			out.Close()

			//If we don't have output then we can delete this file, and stop
			// processing this year because we have gone below the
			// threshold of popularity
			if !hasOutput {
				os.Remove(outFilePath)
				continue
			}
		}

	}
}

func outFileName(file os.FileInfo, fileCount int) string {
	return strings.Split(file.Name(), ".")[0] + "-" + strconv.Itoa(fileCount) + ".txt"
}

//manipulateLine adds the rank and year to the name.
// returns the manipulated line and a boolean value saying if the line should be used.
// Only lines with values greater than 100 incidents should be used.
func manipulateLine(line []string, fileName string, count *yearCount) ([]string, bool) {
	occurances, _ := strconv.Atoi(line[2])
	if occurances > 100 {
		if line[1] == "M" {
			if occurances != count.mOccur {
				count.mRank++
				count.mOccur = occurances
			}
			line = append(line, strconv.Itoa(count.mRank))
		} else {
			if occurances != count.fOccur {
				count.fRank++
				count.fOccur = occurances
			}
			line = append(line, strconv.Itoa(count.fRank))
		}
		line = append(line, extractYear(fileName))
		return line, true
	}
	return nil, false
}

func setLoopLimit(x, chunkSize int, linesArray [][]string) int {
	var limit int
	if x+chunkSize < len(linesArray) {
		limit = x + chunkSize
	} else {
		limit = len(linesArray)
	}
	return limit
}

//Extracts the year from a file name in the format of yob1984.  Returns 1984
func extractYear(fileName string) string {
	return fileName[3:7]
}
