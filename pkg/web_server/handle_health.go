package web_server

import "net/http"

func (s *WebServer) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.responseWithSuccessJSON(w, true)
	}
}

