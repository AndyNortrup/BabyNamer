package usage

import (
	"sort"

	"golang.org/x/net/context"
)

type UsageList []*SettingUsage

func (list UsageList) Len() int {
	return len(list)
}

func (list UsageList) Less(i, j int) bool {
	return list[i].Plain < list[j].Plain
}

func (list UsageList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func NewUsageList(ctx context.Context, user string) UsageList {

	userUsages := getAllUserUsages(user, ctx)
	allUsages := GetNameOrigins()

	return combineUsageLists(userUsages, allUsages, user)
}

func combineUsageLists(
	userUsages map[string]*SettingUsage,
	allUsages map[string]NameOrigin,
	email string) UsageList {

	output := UsageList{}

	for k := range allUsages {
		if userUsages[k] != nil {
			output = append(output, userUsages[k])
		} else {
			output = append(output,
				&SettingUsage{
					NameOrigin: allUsages[k],
					Enabled:    false,
					User:       email,
				})
		}
	}

	sort.Sort(output)
	return output
}
