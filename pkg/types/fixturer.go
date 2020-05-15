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
		Text:      f.Markdown("single-text"),
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
		Text:         f.Markdown("group-text"),
		Published:    true,
		IsSingle:     false,
	}
}

func (f *Fixturer) String(prefix string) string {
	return fmt.Sprintf(`%s-%s`, f.name, prefix)
}

func (f *Fixturer) Markdown(prefix string) string {
	return fmt.Sprintf("**%s**-*%s*\n\n%s", f.name, prefix, `
Lorem *ipsum* dolor sit amet, _consectetur_ adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. 

Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
`)
}
