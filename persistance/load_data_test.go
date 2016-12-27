package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"os"
	"testing"
)

var lineStrings [][]string = [][]string{{"Mary", "F", "7065"},
	{"Anna", "F", "2604"},
	{"Pat", "M", "123"}}

func TestReadLines(t *testing.T) {

	testFilePath := "persistance/names/yob1880.txt"
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("action=TestReadLines error=%v", err.Error())
	}
	out, err := readLines(file)
	for key, outVal := range out {
		for index, value := range outVal {
			if value != lineStrings[key][index] {
				t.Logf("action=TestReadLines result=failed expected=%v recieved=%v",
					lineStrings[key][index], value)
				t.Fail()
			}
		}
	}

}

func TestConvertLineToStat(t *testing.T) {
	out := make(chan *names.Name)
	in := make(chan *names.Name)

	go convertLinesToStat(lineStrings, 1880, out)
	go nameValues(in)
	for x := 0; x < 3; x++ {
		compare := <-in
		name := <-out
		if !compareNames(name, compare) {
			t.Fail()
			t.Logf("action=TestConvertLineToStat \nexpected=%#v \nrecieved=%#v", compare, name)
		}
	}
}

func nameValues(out chan<- *names.Name) {
	defer close(out)

	name := names.NewName("Mary", names.Female)
	name.AddStat(names.NewNameStat(1880, 1, 7065))
	out <- name

	name = names.NewName("Anna", names.Female)
	name.AddStat(names.NewNameStat(1880, 2, 2604))
	out <- name

	name = names.NewName("Pat", names.Male)
	name.AddStat(names.NewNameStat(1880, 1, 123))
	out <- name
}

func compareNames(one, two *names.Name) bool {
	if one.Name != two.Name &&
		one.Gender != two.Gender &&
		one.Stats[1880].Occurrences != two.Stats[1880].Occurrences &&
		one.Stats[1880].Rank != two.Stats[1880].Rank {

		return false
	}
	return true
}
