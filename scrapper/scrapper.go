package scrapper

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/kkucherenkov/golang-test-task/repository"
)

// Scrap make actual request to the site
func Scrap(ur *url.URL) (latency int64, status int, err error) {
	var uri string
	if strings.HasPrefix(ur.String(), "http") {
		uri = ur.String()
	} else {
		uri = "http://" + ur.String()
	}
	t := time.Now().Unix()
	res, err := http.Head(uri)
	t = time.Now().Unix() - t
	if err != nil {
		return t, -1, err
	}
	return t, res.StatusCode, err

}

// Service interface to scrapper process
type Service interface {
	StartScrapWithContext(ctx context.Context) error
}

type service struct {
	repo repository.Repository
	// scrapper Scrapper
}

// NewHTTPService creates new http service
func NewHTTPService(repo repository.Repository) Service {
	s := &service{repo: repo}
	// s.scrapper = newHTTPScrapper()
	return s
}

func (s *service) StartScrapWithContext(ctx context.Context) (err error) {
	var (
		hosts = s.repo.GetAllHosts()
	)
	done := make(chan struct{})
	defer close(done)

	c, errc := s.scrapAllSites(ctx, done, hosts)

	for r := range c {
		if r.err != nil {
			return r.err
		}

		s.repo.Store(r.domain, r.latency, r.status)
	}
	if err := <-errc; err != nil {
		return err
	}
	return nil
}

type result struct {
	domain  string
	latency int64
	status  int
	err     error
}

func (s *service) scrapAllSites(ctx context.Context, done <-chan struct{}, hosts []string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		for _, host := range hosts {
			wg.Add(1)
			go func(domain string) {
				domain, latency, status, err := s.scrapSite(domain)
				if err != nil {
					errc <- err
				}
				select {
				case c <- result{domain, latency, status, err}:
				case <-done:
				}
				wg.Done()
			}(host)
		}

		go func() {
			wg.Wait()
			close(c)
		}()

	}()
	return c, errc
}

func (s *service) scrapSite(site string) (domain string, latency int64, status int, err error) {
	d, err := s.repo.GetByDomain(site)
	// If domain not found must be exist a function
	if err != nil {
		return
	}
	// get first scrapper by cursor
	l, ss, err := Scrap(d.Domain)
	if err != nil {
		ss = 0
	}
	return d.Domain.String(), l, ss, err
}
