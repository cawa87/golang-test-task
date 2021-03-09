package controller

import "context"

// ScrapperController dfdgd
type ScrapperController interface {
	GetMaxLatency(ctx context.Context) (domain string, latency int64, code int, err error)
	GetMinLatency(ctx context.Context) (domain string, latency int64, code int, err error)
	GetByDomain(ctx context.Context, domain string) (latency int64, code int, err error)
	GetCount(ctx context.Context) int
}
