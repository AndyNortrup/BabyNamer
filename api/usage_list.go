package main

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
