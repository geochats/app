package web_server

import "net/http"

func (s *WebServer) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stat, err := s.store.Ping()
		if err == nil {
			s.responseWithSuccessJSON(w, stat)
		} else {
			s.responseWithErrorJSON(w, err)
		}
	}
}

