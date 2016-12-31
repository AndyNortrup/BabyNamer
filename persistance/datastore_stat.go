package persist

import (
	"github.com/AndyNortrup/baby-namer/names"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"strconv"
)

type datastoreStat struct {
	NameKey string
	names.Stat
}

func NewDatastoreStat(name *names.Name, stat names.Stat) *datastoreStat {
	return &datastoreStat{NameKey: name.Key(), Stat: stat}
}

func (stat *datastoreStat) newStatKey(ctx context.Context) *datastore.Key {
	yearStr := strconv.Itoa(stat.Year)
	return datastore.NewKey(ctx, entityTypeStats, stat.NameKey+"::"+yearStr, 0, nil)
}
