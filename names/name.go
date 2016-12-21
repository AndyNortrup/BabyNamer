package names

import (
	"time"
)

type Name struct {
	Name   string
	Gender string
	Stats  map[int]*Stat `datastore:"-"`
}

func NewName(name, gender string) *Name {
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

func (name *Name) MostPopularYear() *Stat {
	year := 0
	for _, stat := range name.Stats {
		if year == 0 {
			year = stat.Year
		} else if stat.Occurrences > name.Stats[year].Occurrences {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) Key() string {
	return name.Name + "::" + name.Gender
}
