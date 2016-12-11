package ssa_data

type SSANameStat struct {
	Occurrences int
	Year        int
	Rank        int
}

func NewSSANameStat(year, rank, occurances int) *SSANameStat {
	return &SSANameStat{
		Year:        year,
		Rank:        rank,
		Occurrences: occurances,
	}
}
