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
	ChatID    int64
	Name      string
	Username  string
	Latitude  float64
	Longitude float64
	Text      string
	Published bool
}

func (p *Point) PublicID() string {
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

type Image struct {
	Width  int32
	Height int32
	Path   string
}
