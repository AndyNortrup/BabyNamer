package usage

import (
	"golang.org/x/net/context"
	"math/rand"
	"time"
)

type UsageGenerator struct {
	ctx  context.Context
	user string
}

func NewUsageGenerator(ctx context.Context, user string) *UsageGenerator {
	return &UsageGenerator{
		ctx:  ctx,
		user: user,
	}
}

func (gen *UsageGenerator) RandomUsage() string {
	usages := []string{"eng", "iri", "sco", "fre", "wel"}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	return usages[r.Intn(len(usages)-1)]
}
