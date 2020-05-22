package web_server

import (
	"fmt"
	"geochats/pkg/markdown"
	"geochats/pkg/storage"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type WebServer struct {
	addr       string
	docRootDir string
	store      *storage.Storage
	logger     *logrus.Entry
}

func New(addr string, documentRootDir string, store *storage.Storage, logger *logrus.Logger) *WebServer {
	return &WebServer{
		addr:       addr,
		docRootDir: documentRootDir,
		store:      store,
		logger:     logger.WithField("package", "web_server"),
	}
}

func (s *WebServer) router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/list.json", s.handleListJSON()).Methods("GET")
	r.HandleFunc("/list.html", s.handleListHTML()).Methods("GET")
	r.HandleFunc("/points/{hashID}/{title}.html", s.handlePointHTML()).Methods("GET")
	r.HandleFunc("/health", s.handleHealth()).Methods("GET")
	r.HandleFunc("/sitemap.xml", s.handleSitemap()).Methods("GET")
	r.HandleFunc("/", s.handleIndex()).Methods("GET")
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir(s.docRootDir))).Methods("GET")
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

func (s *WebServer) templateFunctions() map[string]interface{} {
	functions := sprig.GenericFuncMap()

	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		logrus.Panic(err)
	}
	functions["human_url"] = func(text string) string {
		if len(text) > 40 {
			text = text[0:40]
		}
		text = reg.ReplaceAllString(text, "")
		text = strings.ReplaceAll(text, "__", "_")
		return url.QueryEscape(text)
	}
	functions["md2html"] = func(md string) template.HTML {
		return template.HTML(markdown.ToHTML(md))
	}
	return functions
}

func (s *WebServer) parseTemplate(templateFiles ...string) *template.Template {
	return template.Must(
		template.New("base").
			Funcs(s.templateFunctions()).
			ParseFiles(templateFiles...))
}