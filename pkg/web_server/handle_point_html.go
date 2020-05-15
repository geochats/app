package web_server

import (
	"geochats/pkg/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func (s *WebServer) handlePointHTML() http.HandlerFunc {
	ts := s.parseTemplate(
		"./pkg/web_server/templates/point.go.html",
		"./pkg/web_server/templates/layout.go.html",
	)
	return func(w http.ResponseWriter, r *http.Request) {
		hashID := mux.Vars(r)["hashID"]
		if hashID == "" {
			http.Error(w, "Not found", 404)
			return
		}

		point := new(types.Point)
		if os.Getenv("RANDOM_LIST") != "" {
			f := types.NewRandomFixturer("fake")
			p := f.Group()
			point = &p
		} else {
			points, err := s.store.ListPublishedPoints()
			if err != nil {
				logrus.Error(err.Error())
				http.Error(w, "Internal Server Error", 500)
			}
			for _, p := range points {
				if p.HashedID() == hashID {
					point = &p
				}
			}
		}

		if point == nil {
			http.Error(w, "Not found", 404)
			return
		}

		if err := ts.Execute(w, point); err != nil {
			logrus.Error(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
