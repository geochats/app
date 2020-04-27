package web_server

import "net/http"

func (s *WebServer) routes() {
	s.router.HandleFunc("/list", s.handleList()).Methods("GET")
	s.router.HandleFunc("/health", s.handleHealth()).Methods("GET")
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public"))).Methods("GET")
}
