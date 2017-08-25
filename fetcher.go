package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"
)

// DocumentFetcher provides async interface to fetch bulk of urls.
// Returned channel is unbuffered and will be closed when all documents will are fetched.
type DocumentFetcher interface {
	Fetch(urls ...url.URL) (<-chan RawDocument, error)
}

func NewDocumentFetcher(cfg FetcherConfig, l *logrus.Logger) (DocumentFetcher, error) {
	if cfg.WorkerCount <= 0 {
		return nil, errors.New("invalid number of workers")
	}

	workers := make(chan struct{}, cfg.WorkerCount)
	for i := 0; i < cfg.WorkerCount; i++ {
		workers <- struct{}{}
	}

	httpClient := &http.Client{
		Transport: &http.Transport{},
		Timeout:   cfg.Timeout,
	}

	fetcher := &documentFetcher{
		workers: workers,

		http: httpClient,
		log:  l.WithField(logPlace, "FETCHER"),
	}

	return fetcher, nil
}

// documentFetcher implements DocumentParser by limting concurrency of document
// fetching (if HTML processing is the bottleneck there will be no value in simultaneous
// fetching of all documents - the queues will build up, huge number of goroutines wiil
// consume RAM and and scheduler resources)
type documentFetcher struct {
	workers chan struct{}

	http *http.Client
	log  *logrus.Entry
}

func (f *documentFetcher) Fetch(urls ...url.URL) (<-chan RawDocument, error) {
	rawDocs := make(chan RawDocument)

	go func() {
		wg := sync.WaitGroup{}

		for i := range urls {
			// Wait for available worker
			<-f.workers
			wg.Add(1)

			go func(url url.URL) {
				// Release worker and signal to caller
				defer func() {
					f.workers <- struct{}{}
					wg.Done()
				}()

				if rd, err := f.fetchDoc(url); err == nil {
					rawDocs <- rd
				} else {
					f.log.WithError(err).WithField("url", url).Error("Failed to fetch document")
				}
			}(urls[i])
		}

		// Wait for all workers and signal to caller
		wg.Wait()
		close(rawDocs)
	}()

	return rawDocs, nil
}

func (f *documentFetcher) fetchDoc(url url.URL) (RawDocument, error) {
	rd := RawDocument{
		Url: url.String(),
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return rd, err
	}

	resp, err := f.http.Do(req)
	if err != nil {
		return rd, err
	}
	defer resp.Body.Close()

	rd.Meta.Status = resp.StatusCode

	if rd.Meta.Status >= 200 && rd.Meta.Status < 300 {
		rd.Meta.ContentType = resp.Header.Get("Content-Type")

		if rd.Meta.ContentType == "text/html" {
			if rd.Content, err = ioutil.ReadAll(resp.Body); err != nil {
				return rd, err
			}
			rd.Meta.ContentLength = len(rd.Content)
		}
	}

	return rd, nil
}
