package ssa_data

import "time"

type Name struct {
	Name   string
	Gender string
	Stats  map[int]*SSANameStat
}

func NewName(name, gender string) *Name {
	result := &Name{
		Name:   name,
		Gender: gender,
	}
	result.Stats = make(map[int]*SSANameStat)
	return result
}

func (name *Name) addStat(stat *SSANameStat) {
	name.Stats[stat.Year] = stat
}

func (name *Name) LatestYear() *SSANameStat {
	year := 0
	for _, stat := range name.Stats {
		if stat.Year > year {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) FirstYear() *SSANameStat {
	year := time.Now().Year()
	for _, stat := range name.Stats {
		if stat.Year < year {
			year = stat.Year
		}
	}
	return name.Stats[year]
}

func (name *Name) MostPopularYear() *SSANameStat {
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
