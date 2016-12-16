package ssa_data

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
