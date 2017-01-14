package names

import (
	"fmt"
	"sort"
	"time"
)

type Name struct {
	Name   string
	Gender Gender
	Stats  map[int]*Stat `datastore:"-"`
}

func NewName(name string, gender Gender) *Name {
	result := &Name{
		Name:   name,
		Gender: gender,
	}
	result.Stats = make(map[int]*Stat)
	return result
}

func (name *Name) AddStat(stat *Stat) {
	if name.Stats == nil {
		name.Stats = make(map[int]*Stat)
	}
	name.Stats[stat.Year] = stat
}

func (name *Name) LatestYear() *Stat {
	year := 0
	for _, stat := range name.Stats {
		if stat.Year > year {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) FirstYear() *Stat {
	year := time.Now().Year()
	for _, stat := range name.Stats {
		if stat.Year < year {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) HighestRank() *Stat {
	var year int
	for _, stat := range name.SortedStats() {
		if year == 0 {
			year = stat.Year
		} else if stat.Rank < name.Stats[year].Rank {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) LowestRank() *Stat {
	var year int
	for _, stat := range name.SortedStats() {
		if year == 0 {
			year = stat.Year
		} else if stat.Rank > name.Stats[year].Rank {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) MostOccurrences() *Stat {
	year := 0
	for _, stat := range name.Stats {
		if year == 0 || stat.Occurrences > name.Stats[year].Occurrences {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) Key() string {
	return name.Name + "::" + name.Gender.GoString()
}

func (name *Name) SortedStats() []*Stat {
	years := []int{}
	result := []*Stat{}

	for _, stat := range name.Stats {
		years = append(years, stat.Year)
	}
	sort.Ints(years)

	for _, year := range years {
		result = append(result, name.Stats[year])
	}

	return result
}

func (name *Name) GoString() string {
	str := fmt.Sprintf("Name=%v Gender=%v ", name.Name, name.Gender)
	for _, value := range name.Stats {
		str = fmt.Sprintf("%v Year=%v Rank=%v Occurance=%v", str, value.Year, value.Rank, value.Occurrences)
	}
	return str
}
