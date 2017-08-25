package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type Server struct {
	fetcher DocumentFetcher
	parser  DocumentParser

	log *logrus.Entry

	http *http.Server
}

func NewServer(cfg ServerConfig, f DocumentFetcher, p DocumentParser, l *logrus.Logger) (*Server, error) {
	server := &Server{
		fetcher: f,
		parser:  p,

		log: l.WithField(logPlace, "HTTP"),

		http: &http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/scrap", server.requestHandler)
	server.http.Handler = mux

	return server, nil
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.http.Shutdown(context.Background())
}

func (s *Server) requestHandler(w http.ResponseWriter, req *http.Request) {
	l := s.log.WithFields(logrus.Fields{
		"path":   req.URL.Path,
		"method": req.Method,
	})
	l.Info("Incoming request")

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reqData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rawUrls []string
	if err := json.Unmarshal(reqData, &rawUrls); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	urls := make([]url.URL, 0, len(rawUrls))
	for _, r := range rawUrls {
		u, err := url.Parse(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid url - %s", u.String())))
			return
		}
		urls = append(urls, *u)
	}

	rawDocs, err := s.fetcher.Fetch(urls...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	docs, err := s.parser.Parse(rawDocs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respData, err := json.Marshal(&docs)
	if err != nil {
		s.log.WithError(err).Error("Failed to marsahl JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(respData); err != nil {
		s.log.WithError(err).Error("Failed to write response")
	}
}
