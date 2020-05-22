package web_server

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *WebServer) handleIndex() http.HandlerFunc {
	ts := s.parseTemplate(
		"./public/index.go.html",
		"./public/layout.go.html",
	)
	return func(w http.ResponseWriter, r *http.Request) {
		err := ts.Execute(w, nil)
		if err != nil {
			logrus.Error(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
