package types

import (
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type Group struct {
	ChatID           int64
	Title        string
	Username     string
	Userpic      Image
	MembersCount int32
	Latitude     float64
	Longitude    float64
	Description  string
}


func (g *Group) Complete() bool {
	return g.Username != "" && g.Latitude != 0
}

func (g *Group) DescriptionHTML()  string {
	h := string(markdown.ToHTML([]byte(g.Description), nil, nil))
	replacer := strings.NewReplacer(">http://", ">", ">https://", ">")
	h = replacer.Replace(h)
	return bluemonday.UGCPolicy().Sanitize(h)
}