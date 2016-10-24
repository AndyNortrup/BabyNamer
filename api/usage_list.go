package main

import (
	"context"
	"sort"

	"google.golang.org/appengine/user"
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

func usageListFromDatastore(ctx context.Context, user *user.User) UsageList {

	userUsages := getAllUserUsages(user.Email, ctx)
	allUsages := getNameOrigins()

	return combineUsageLists(userUsages, allUsages, user)
}

func combineUsageLists(
	userUsages map[string]*SettingUsage,
	allUsages map[string]NameOrigin,
	user *user.User) UsageList {

	output := UsageList{}

	for k := range allUsages {
		if userUsages[k] != nil {
			output = append(output, userUsages[k])
		} else {
			output = append(output,
				&SettingUsage{
					NameOrigin: allUsages[k],
					Enabled:    false,
					User:       user.Email,
				})
		}
	}

	sort.Sort(output)
	return output
}
