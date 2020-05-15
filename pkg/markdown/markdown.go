package markdown

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

func ToHTML(text string) string {
	flags := html.CommonFlags | html.SkipHTML
	renderer := html.NewRenderer(html.RendererOptions{Flags: flags})
	h := string(markdown.ToHTML([]byte(text), nil, renderer))
	replacer := strings.NewReplacer(">http://", ">", ">https://", ">")
	h = replacer.Replace(h)
	return bluemonday.UGCPolicy().Sanitize(h)
}
