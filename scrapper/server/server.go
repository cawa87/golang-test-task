package server

import (
	"github.com/souz9/golang-test-task/scrapper/hchecker"
	"net/http"
)

type HealthChecker interface {
	Status() []hchecker.Status
	StatusOf(string) hchecker.Status
	FindMin() hchecker.Status
	FindMax() hchecker.Status
}

type Server struct {
	HChecker HealthChecker
}

func (s *Server) router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/status", s.status)
	router.HandleFunc("/status/min", s.statusMin)
	router.HandleFunc("/status/max", s.statusMax)
	return router
}

func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s.router())
}
