package names

type NameList []*Name

func (list NameList) Len() int {
	return len(list)
}

func (list NameList) Less(i, j int) bool {
	return list[i].Name < list[j].Name
}

func (list NameList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
