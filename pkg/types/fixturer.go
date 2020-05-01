package types

import (
	"fmt"
	"math/rand"
	"time"
)

type Fixturer struct {
	name string
}

func NewFixedFixturer(name string) *Fixturer {
	return &Fixturer{name}
}

func NewRandomFixturer(name string) *Fixturer {
	return &Fixturer{fmt.Sprintf("%s-%d", name, time.Now().UnixNano())}
}

func (f *Fixturer) Single() Point {
	return Point{
		ChatID:    int64(rand.Uint64()),
		Username:  f.String("username"),
		Latitude:  50 + 5*rand.NormFloat64(),
		Longitude: 50 + 5*rand.NormFloat64(),
		Text:      f.String("text"),
		Published: true,
		IsSingle:  true,
	}
}

func (f *Fixturer) Group() Point {
	return Point{
		ChatID:       int64(rand.Uint64()),
		Username:     f.String("username"),
		Latitude:     50 + 5*rand.NormFloat64(),
		Longitude:    50 + 5*rand.NormFloat64(),
		MembersCount: int32(rand.Intn(1000)),
		Text:         f.String("desc"),
		Published:    true,
		IsSingle:     false,
	}
}

func (f *Fixturer) String(prefix string) string {
	return fmt.Sprintf(`%s-%s`, f.name, prefix)
}
