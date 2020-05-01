package types

import (
	"crypto/sha256"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/microcosm-cc/bluemonday"
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

func (p *Point) HashedID() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", p.ChatID)))
	return fmt.Sprintf("%x", h)
}

func (p *Point) TextHTML() string {
	flags := html.CommonFlags | html.SkipHTML
	renderer := html.NewRenderer(html.RendererOptions{Flags: flags})
	h := string(markdown.ToHTML([]byte(p.Text), nil, renderer))
	replacer := strings.NewReplacer(">http://", ">", ">https://", ">")
	h = replacer.Replace(h)
	return bluemonday.UGCPolicy().Sanitize(h)
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