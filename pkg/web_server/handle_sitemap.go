package web_server

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"text/template"
)

func (s *WebServer) handleSitemap() http.HandlerFunc {
	ts := template.Must(
		template.New("sitemap.go.xml").
		Funcs(s.templateFunctions()).
		ParseFiles("./public/sitemap.go.xml"))
	return func(w http.ResponseWriter, r *http.Request) {
		points, err := s.store.ListPublishedPoints()
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't get points list: %v", err))
			return
		}
		if os.Getenv("RANDOM_LIST") != "" {
			f := types.NewRandomFixturer("fake")
			for i := 0; i < 100; i++ {
				points = append(points, f.Single())
			}
			for i := 0; i < 10; i++ {
				points = append(points, f.Group())
			}
		}
		w.Header().Set("Content-Type", "application/xml")
		if err := ts.Execute(w, points); err != nil {
			logrus.Error(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
