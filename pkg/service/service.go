package service

import (
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/l-vitaly/golang-test-task/pkg/crawl"
)

// error consts
var (
	ErrEmptyURLs = errors.New("empty urls")
)

type jobResult struct {
	result crawl.Result
	err    error
}

// Service service interface
type Service interface {
	PostURLs(urls []string) ([]crawl.Result, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(crawl crawl.Crawl, maxWorkers int, logger log.Logger) Service {
	var svc Service
	{
		svc = NewBasicService(crawl, maxWorkers)
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewBasicService returns a native, stateless implementation of Service.
func NewBasicService(crawl crawl.Crawl, maxWorkers int) Service {
	return &basicService{
		crawl:      crawl,
		maxWorkers: maxWorkers,
	}
}

type basicService struct {
	crawl      crawl.Crawl
	maxWorkers int
}

func (s *basicService) worker(urls <-chan string, results chan<- jobResult) {
	for url := range urls {
		r, err := s.crawl.Do(url)
		results <- jobResult{
			result: r,
			err:    err,
		}
	}
}

// PostURLs run parse urls.
func (s *basicService) PostURLs(urls []string) ([]crawl.Result, error) {
	if len(urls) == 0 {
		return nil, ErrEmptyURLs
	}
	crawlResult := []crawl.Result{}

	urlLen := len(urls)

	jobs := make(chan string, urlLen)
	jobResults := make(chan jobResult, urlLen)

	workers := urlLen
	if workers > s.maxWorkers {
		workers = s.maxWorkers
	}

	for w := 1; w <= workers; w++ {
		go s.worker(jobs, jobResults)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	for i := 1; i <= urlLen; i++ {
		jobResult := <-jobResults
		if jobResult.err != nil {
			continue
		}
		crawlResult = append(crawlResult, jobResult.result)
	}

	return crawlResult, nil
}
