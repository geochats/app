package web_server

import (
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/loaders"
	"geochats/pkg/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type WebServer struct {
	addr       string
	docRootDir string
	tg         client.AbstractClient
	store      *storage.Storage
	loader     *loaders.ChannelInfoLoader
	logger     *logrus.Entry
}

func New(addr string, documentRootDir string, tgClient client.AbstractClient, store *storage.Storage, loader *loaders.ChannelInfoLoader, logger *logrus.Logger) *WebServer {
	return &WebServer{
		addr:       addr,
		docRootDir: documentRootDir,
		tg:         tgClient,
		store:      store,
		loader:     loader,
		logger:     logger.WithField("package", "web_server"),
	}
}

func (s *WebServer) router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/list", s.handleList()).Methods("GET")
	r.HandleFunc("/health", s.handleHealth()).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(s.docRootDir))).Methods("GET")
	return r
}

func (s *WebServer) Listen() error {
	srv := &http.Server{
		Handler:      handlers.LoggingHandler(s.logger.Writer(), s.router()),
		Addr:         s.addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	s.logger.Infof("Listen on %s", s.addr)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("can't start web server: %v", err)
	}
	return nil
}
