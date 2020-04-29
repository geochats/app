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

func (f *Fixturer) Point() Point {
	return Point{
		ChatID:    int64(rand.Uint64()),
		Name:      f.String("name"),
		Username:  f.String("username"),
		Latitude:  50 + 5*rand.NormFloat64(),
		Longitude: 50 + 5*rand.NormFloat64(),
		Text:      f.String("text"),
		Published: true,
	}
}

func (f *Fixturer) Group() Group {
	return Group{
		ChatID:   int64(rand.Uint64()),
		Title:    f.String("title"),
		Username: f.String("username"),
		Userpic: Image{
			Width:  1024,
			Height: 1024,
			Path:   "https://picsum.photos/1024/1024",
		},
		MembersCount: int32(rand.Intn(1000)),
		Latitude:     50 + 5*rand.NormFloat64(),
		Longitude:    50 + 5*rand.NormFloat64(),
		Text:         f.String("desc"),
		Published:    true,
	}
}

func (f *Fixturer) String(prefix string) string {
	return fmt.Sprintf(`%s-%s`, f.name, prefix)
}
