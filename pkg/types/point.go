package types

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Point struct {
	ChatID       int64
	Username     string
	Text         string
	Latitude     float64
	Longitude    float64
	MembersCount int32
	Published    bool
	IsSingle     bool
}

func (p *Point) Title() string {
	title := fmt.Sprintf("%.80s", p.Text)
	title = strings.ReplaceAll(title, "/", "_")
	title = strings.ReplaceAll(title, " ", "_")
	return title
}

func (p *Point) HashedID() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", p.ChatID)))
	return fmt.Sprintf("%x", h)
}

func (p *Point) PublicURI() string {
	flag := "g"
	id := fmt.Sprintf("%d", p.ChatID)
	if p.IsSingle {
		flag = "s"
		id = p.HashedID()
	}
	return fmt.Sprintf("https://miting.link/#%s:%s", flag, id)
}