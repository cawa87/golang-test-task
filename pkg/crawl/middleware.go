package crawl

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(Crawl) Crawl

type loggingMiddleware struct {
	logger log.Logger
	next   Crawl
}

func (m loggingMiddleware) Do(url string) (result Result, err error) {
	defer func(begin time.Time) {
		m.logger.Log("method", "Do", "url", url, "took", time.Since(begin), "err", err)
	}(time.Now())
	return m.next.Do(url)
}

// LoggingMiddleware takes a logger as a dependency
// and returns a LoggingMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Crawl) Crawl {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}
