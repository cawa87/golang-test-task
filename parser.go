package main

import (
	"bytes"
	"errors"
	"sync"

	"golang.org/x/net/html"

	"github.com/sirupsen/logrus"
)

// DocumentParser provider interface to parse HTML documents.
// Channel, passed to Parse must be unbuffered and must be closed by user.
type DocumentParser interface {
	Parse(rawDocs <-chan RawDocument) ([]Document, error)
}

func NewDocumentParser(cfg ParserConfig, l *logrus.Logger) (DocumentParser, error) {
	if cfg.WorkerCount <= 0 {
		return nil, errors.New("invalid number of workers")
	}

	workers := make(chan struct{}, cfg.WorkerCount)
	for i := 0; i < cfg.WorkerCount; i++ {
		workers <- struct{}{}
	}

	parser := &documentParser{
		workers: workers,

		log: l.WithField(logPlace, "PARSER"),
	}

	return parser, nil
}

type documentParser struct {
	workers chan struct{}

	log *logrus.Entry
}

// documentParser implements DocumentParser by limting concurrency of HTML
// parsing (huge number fo goroutines does not speed up completion of CPU-intensive task,
// by consumes RAM and scheduler resources)
func (p *documentParser) Parse(rawDocs <-chan RawDocument) ([]Document, error) {
	var (
		docs = make([]Document, 0, len(rawDocs))
		lock sync.Mutex
	)

	wg := sync.WaitGroup{}

	for r := range rawDocs {
		// Wait for available worker
		<-p.workers
		wg.Add(1)

		go func(raw RawDocument) {
			// Release worker and signal to caller
			defer func() {
				p.workers <- struct{}{}
				wg.Done()
			}()

			if d, err := p.parseDoc(raw); err == nil {
				lock.Lock()
				docs = append(docs, d)
				lock.Unlock()
			} else {
				p.log.WithError(err).WithField("url", raw.Url).Error("Failed to parse document")
			}
		}(r)
	}

	wg.Wait()

	return docs, nil
}

func (p *documentParser) parseDoc(raw RawDocument) (Document, error) {
	doc := Document{
		Url:  raw.Url,
		Meta: raw.Meta,
	}

	node, err := html.Parse(bytes.NewReader(raw.Content))
	if err != nil {
		return doc, err
	}

	var (
		tags     = make(map[string]int, 0)
		traverse func(*html.Node)
	)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			tags[n.Data] = tags[n.Data] + 1
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(node)

	doc.Elements = make([]DocumentElement, 0, len(tags))
	for t, c := range tags {
		e := DocumentElement{
			TagName: t,
			Count:   c,
		}
		doc.Elements = append(doc.Elements, e)
	}

	return doc, nil
}
