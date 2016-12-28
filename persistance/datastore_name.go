package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
)

type datastoreName struct {
	names.Name
	Gender string
	Random float32
}

func newDatastoreName(name *names.Name) *datastoreName {
	dName := &datastoreName{
		Name:   *name,
		Gender: name.Gender.GoString(),
		Random: randomFloat(),
	}
	return dName
}
