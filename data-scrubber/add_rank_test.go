package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestProcessFolder(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	inFiles := path.Join(pwd, "data-scrubber", "inFiles")
	outFiles := path.Join(pwd, "data-scrubber", "outFiles")
	files, err := ioutil.ReadDir(inFiles)
	if err != nil {
		log.Printf("Unable to read files in directory. Error=%v", err.Error())
		t.FailNow()
	}
	processFolder(files, inFiles, outFiles, 100)
}

func TestManipulateLine(t *testing.T) {
	fileName := "yob1880.txt"
	count := &yearCount{fRank: 1, mRank: 1, fOccur: 102, mOccur: 102}

	inputLines := [][]string{
		{"Sheppard", "M", "5"},
		{"Sheppard", "M", "100"},
		{"Sheppard", "M", "101"},
		{"Sheppard", "M", "99"},
	}
	outputLines := [][]string{
		{"Sheppard", "M", "5"},
		{"Sheppard", "M", "100"},
		{"Sheppard", "M", "101", "2", "1880"},
		{"Sheppard", "M", "99"},
	}

	useValues := []bool{false, false, true, false}

	for index, line := range inputLines {
		outLine, use := manipulateLine(line, fileName, count)
		if useValues[index] != use {
			t.Log("Wrong use value returned. Expected=%v Recieved=%v for line=%v",
				useValues[index], use, line)
			t.Fail()
		}
		for key, value := range outLine {
			if outputLines[index][key] != value {
				t.Log("Wrong line returned. expected=%v recieved=%v",
					outputLines[index], outLine)
			}
		}
	}
}
