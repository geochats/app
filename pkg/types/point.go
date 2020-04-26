package types

import (
	"crypto/sha256"
	"fmt"
)

type Point struct {
	ChatID      string     `json:"id"`
	Username    string    `json:"username"`
	Photo       Image     `json:"photo"`
	Coords      []float64 `json:"coords"`
	Description string    `json:"description"`
}

func NewPoint(chatId int64) *Point {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", chatId)))
	return &Point{
		ChatID: fmt.Sprintf("%x", h),
	}
}

func (p *Point) Complete() bool {
	return p.Photo.Path != "" && len(p.Coords) > 0 && p.Description != ""
}

type Image struct {
	Width    int32
	Height   int32
	Path     string
}
