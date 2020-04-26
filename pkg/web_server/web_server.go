package web_server

import (
	"encoding/json"
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"geochats/pkg/types"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type WebServer struct {
	addr   string
	tg     client.AbstractClient
	store storage.Storage
	loader *loaders.ChannelInfoLoader
	router *mux.Router
	logger *logrus.Entry
}

func New(addr string, tgClient client.AbstractClient, store storage.Storage, loader *loaders.ChannelInfoLoader, logger *logrus.Logger) *WebServer {
	return &WebServer{
		addr:   addr,
		tg:     tgClient,
		store: store,
		loader: loader,
		router: mux.NewRouter(),
		logger: logger.WithField("package", "web_server"),
	}
}

func (s *WebServer) Listen() error {
	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, s.router),
		Addr:         s.addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	s.routes()
	s.logger.Infof("Listen on %s", s.addr)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("can't start web server: %v", err)
	}
	return nil
}

func (s *WebServer) routes() {
	s.router.HandleFunc("/add", s.handleAdd()).Methods("POST")
	s.router.HandleFunc("/list", s.handleList()).Methods("GET")
	s.router.HandleFunc("/health", s.handleHealth()).Methods("GET")
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public"))).Methods("GET")
}

func (s *WebServer) handleAdd() http.HandlerFunc {
	type reqSpec struct {
		Coords []float64 `json:"coords"`
		Link   string    `json:"link"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req reqSpec
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("cant' decode add request: %v", err))
			return
		}
		defer func() { _ = r.Body.Close() }()

		username := strings.Replace(req.Link, "https://t.me/", "", -1)
		gr, err := s.loader.Export(username)
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't load tg info about group: %v", err))
			return
		}

		gr.Coords = req.Coords
		if err := s.store.AddGroup(gr); err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't store new group: %v", err))
			return
		}

		s.responseWithSuccessJSON(w, true)
	}
}

func (s *WebServer) handleList() http.HandlerFunc {
	type resp struct {
		Groups []types.Group
		Points []types.Point
	}
	return func(w http.ResponseWriter, r *http.Request) {
		points, err := s.store.ListPoint()
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't load points: %v", err))
			return
		}
		completePoint := make([]types.Point, 0)
		for _, p := range points {
			s.logger.Infof("point: %#v", p)
			if p.Complete() {
				completePoint = append(completePoint, p)
			}
		}

		groups, err := s.store.ListGroups()
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't load groups: %v", err))
			return
		}

		s.responseWithSuccessJSON(w, resp{
			Groups: groups,
			Points: completePoint,
		})
	}
}

func (s *WebServer) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.responseWithSuccessJSON(w, true)
	}
}
