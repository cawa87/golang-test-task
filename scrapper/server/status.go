package server

import (
	"encoding/json"
	"fmt"
	"github.com/souz9/golang-test-task/scrapper/hchecker"
	"net/http"
)

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if site := q.Get("site"); site != "" {
		s.statusOf(w, site)
		return
	}
	s.writeStatus(w, s.HChecker.Status())
}

func (s *Server) statusOf(w http.ResponseWriter, site string) {
	status := s.HChecker.StatusOf(site)
	if status.Site == "" {
		s.writeStatus(w, nil)
		return
	}
	s.writeStatus(w, []hchecker.Status{status})
}

func (s *Server) statusMin(w http.ResponseWriter, r *http.Request) {
	status := s.HChecker.FindMin()
	if status.Site == "" {
		s.writeStatus(w, nil)
		return
	}
	s.writeStatus(w, []hchecker.Status{status})
}

func (s *Server) statusMax(w http.ResponseWriter, r *http.Request) {
	status := s.HChecker.FindMax()
	if status.Site == "" {
		s.writeStatus(w, nil)
		return
	}
	s.writeStatus(w, []hchecker.Status{status})
}

func (s *Server) writeStatus(w http.ResponseWriter, status []hchecker.Status) {
	if status == nil {
		status = []hchecker.Status{}
	}

	body, err := json.Marshal(status)
	if err != nil {
		http.Error(w, fmt.Sprintf("marhal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
