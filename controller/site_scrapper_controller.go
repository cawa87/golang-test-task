package controller

import (
	"context"

	"github.com/kkucherenkov/golang-test-task/repository"
)

type siteScrapperController struct {
	repository repository.Repository
}

// New creates new site scrapper controller
func New(repo repository.Repository) ScrapperController {
	ssc := &siteScrapperController{}

	ssc.repository = repo
	return ssc
}

func (sc *siteScrapperController) GetMaxLatency(ctx context.Context) (domain string, latency int64, code int, err error) {
	host, er := sc.repository.GetMaxLatency()
	if er != nil {
		return "", 0, 0, er
	}
	return host.Domain.String(), host.Latency, host.HTTPStatusCode, nil
}

func (sc *siteScrapperController) GetMinLatency(ctx context.Context) (domain string, latency int64, code int, err error) {
	host, er := sc.repository.GetMinLatency()
	if er != nil {
		return "", 0, 0, er
	}
	return host.Domain.String(), host.Latency, host.HTTPStatusCode, nil
}
func (sc *siteScrapperController) GetByDomain(ctx context.Context, domain string) (latency int64, code int, err error) {
	host, er := sc.repository.GetByDomain(domain)
	if er != nil {
		return 0, 0, er
	}
	return host.Latency, host.HTTPStatusCode, nil
}
func (sc *siteScrapperController) GetCount(ctx context.Context) int {
	length := sc.repository.GetCount()

	return length
}
