package repository

import (
	"github.com/kkucherenkov/golang-test-task/model"
)

// Repository it's an interface to datasource
type Repository interface {
	Store(domain string, latency int64, code int) error
	GetMaxLatency() (result model.Host, err error)
	GetMinLatency() (result model.Host, err error)
	GetByDomain(domain string) (host model.Host, err error)
	CountSites() int
	GetAllSites() []string
}
