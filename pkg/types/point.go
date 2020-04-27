package types

import (
	"crypto/sha256"
	"fmt"
)

type Point struct {
	ChatID    int64
	Photo     Image
	Latitude  float64
	Longitude float64
	MottoID   string
}

func (p *Point) PublicID() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", p.ChatID)))
	return fmt.Sprintf("%x", h)
}

func (p *Point) Complete() bool {
	return p.Latitude != 0
}

type Image struct {
	Width  int32
	Height int32
	Path   string
}
