package names

import "time"

type Stat struct {
	Occurrences int
	Year        int
	Rank        int
}

func NewNameStat(year, rank, occurances int) *Stat {
	return &Stat{
		Year:        year,
		Rank:        rank,
		Occurrences: occurances,
	}
}

func (stat *Stat) YearAsTime() time.Time {
	return time.Date(stat.Year, 1, 1, 0, 0, 0, 0, time.Local)
}
