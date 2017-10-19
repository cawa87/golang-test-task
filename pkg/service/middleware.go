package service

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/l-vitaly/golang-test-task/pkg/crawl"
)

// Middleware describes a service middleware.
type Middleware func(Service) Service

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (m loggingMiddleware) PostURLs(urls []string) (results []crawl.Result, err error) {
	defer func(begin time.Time) {
		m.logger.Log("method", "PostURLs", "urls", len(urls), "err", err, "took", time.Since(begin))
	}(time.Now())
	return m.next.PostURLs(urls)
}

// LoggingMiddleware takes a logger as a dependency
// and returns a LoggingMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}
