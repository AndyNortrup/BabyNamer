package persist

import "github.com/AndyNortrup/baby-namer/names"

type datastoreName struct {
	names.Name
	Random float32
}

func newDatastoreName(name *names.Name) *datastoreName {
	dName := &datastoreName{
		Name:   *name,
		Random: randomFloat(),
	}
	return dName
}
