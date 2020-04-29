package types

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type Group struct {
	ChatID       int64
	Title        string
	Username     string
	Text         string
	Userpic      Image
	MembersCount int32
	Latitude     float64
	Longitude    float64
	Published    bool
}

func (g *Group) Complete() bool {
	return g.Username != "" && g.Latitude != 0
}

func (g *Group) TextHTML() string {
	flags := html.CommonFlags | html.SkipHTML
	renderer := html.NewRenderer(html.RendererOptions{Flags: flags})
	h := string(markdown.ToHTML([]byte(g.Text), nil, renderer))
	replacer := strings.NewReplacer(">http://", ">", ">https://", ">")
	h = replacer.Replace(h)
	return bluemonday.UGCPolicy().Sanitize(h)
}
