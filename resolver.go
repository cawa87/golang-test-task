package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	xhtml "golang.org/x/net/html"
)

type httpStatus int

const (
	//HTTPBadRequest HTTPBadRequest
	httpBadRequest httpStatus = 400
)

type urlRequestStatus int

const (
	badURL         urlRequestStatus = -1
	errorOnRequest urlRequestStatus = -2
)

type tagCounter struct {
	TagName string `json:"tag-name"`
	Count   int64  `json:"count"`
}
type OneUrlContext struct {
	URL      string
	response []byte
	Meta     struct {
		Status        int    `json:"status"`
		ContentType   string `json:"content-type,omitempty"`
		ContentLength int64  `json:"content-length,omitempty"`
	} `json:"meta"`
	tags     map[string]int64
	Elements []*tagCounter `json:"elements,omitempty"`
}

//HandleContext HandleContext
type HandleContext struct {
	request       []byte
	urlsToResolve []string
	urlContext    []*OneUrlContext
	answerChannel chan int
	answeredCount int
	log           *Logger
}

func (c *HandleContext) prepareURLContextes() {

}

func (c *HandleContext) clear() {
	c.request = c.request[:cap(c.request)]
	c.urlsToResolve = c.urlsToResolve[:0]
	c.log.Debug("len c.urlsToResolve on clear:", len(c.urlsToResolve))
	for _, uc := range c.urlContext {
		uc.response = uc.response[0:cap(uc.response)]
		uc.Elements = uc.Elements[:0]
		for k := range uc.tags {
			delete(uc.tags, k)
		}
	}
	c.urlContext = c.urlContext[:cap(c.urlContext)]

	c.answeredCount = 0
}
func NewHandleContext(log *Logger) *HandleContext {
	c := &HandleContext{
		request:       make([]byte, 4096, 4096),
		urlsToResolve: make([]string, 0, 128),
		urlContext:    make([]*OneUrlContext, 128, 128),
		answerChannel: make(chan int, 128),
		answeredCount: 0,
		log:           log,
	}
	for i := 0; i < len(c.urlContext); i++ {
		c.urlContext[i] = &OneUrlContext{
			response: make([]byte, 32767, 32767),
			tags:     make(map[string]int64),
			Elements: make([]*tagCounter, 0, 128),
		}
	}
	c.clear()
	return c
}

//Resolver Resolver
type Resolver struct {
	waitResponseTimeoutInMs time.Duration
	limitChannel            chan bool
	contextes               []*HandleContext
	contextBusy             []*int32
	countRequests           int64
	jobsChannel             chan func()
	log                     *Logger
}

func NewResolver(waitResponseTimeoutInMs time.Duration, limitRequests int, l *Logger) (r *Resolver) {
	r = &Resolver{
		waitResponseTimeoutInMs: waitResponseTimeoutInMs,
		limitChannel:            make(chan bool, limitRequests),
		contextes:               make([]*HandleContext, limitRequests, limitRequests),
		contextBusy:             make([]*int32, limitRequests, limitRequests),
		log:                     l,
	}
	for _, v := range r.contextBusy {
		var i int32
		v = &i
	}
	for i := 0; i < limitRequests; i++ {
		r.limitChannel <- true
	}
	return r
}
func countTags(context *OneUrlContext) error {

	tokenizer := xhtml.NewTokenizer(bytes.NewReader(context.response))
	var err error

	for err == nil {
		tt := tokenizer.Next()
		switch {
		case tt == xhtml.ErrorToken:
			// End of the document, we're done
			return err
		case tt == xhtml.StartTagToken ||
			tt == xhtml.SelfClosingTagToken:
			t := tokenizer.Token()
			context.tags[t.Data]++
		}
	}
	return err
}

func (resolver *Resolver) resolveOneURL(urlContext *OneUrlContext, timeoutInMs time.Duration) error {
	log := resolver.log
	answer := urlContext

	log.Trace("Start handle url:", urlContext.URL)
	if _, errParse := neturl.Parse(urlContext.URL); errParse != nil {
		log.Trace("Error url:", urlContext.URL)
		answer.Meta.Status = int(badURL)
		return errParse
	}
	client := http.Client{Timeout: timeoutInMs * time.Millisecond}
	log.Trace("Get url:", urlContext.URL)
	response, err := client.Get(urlContext.URL)

	if err != nil {
		log.Trace("Get url error:", urlContext.URL, err)
		answer.Meta.Status = int(errorOnRequest)
		return err
	}
	meta := &answer.Meta
	meta.Status = response.StatusCode
	meta.ContentType = response.Header.Get("Content-Type")
	searchText := "text/html"

	needParseHTML :=
		meta.Status >= 200 &&
			meta.Status < 300 &&
			meta.ContentType[0:len(searchText)] == searchText
	log.Trace("needParseHTML:", needParseHTML, urlContext.URL, meta.ContentType)
	if needParseHTML {

		n, errRead := response.Body.Read(urlContext.response)
		response.Body.Close()
		if errRead != nil && errRead != io.EOF {
			log.Debug("errRead", errRead)
			return errRead
		}
		urlContext.response = urlContext.response[0:n]
		meta.ContentLength = int64(len(urlContext.response))
		log.Debug("meta.ContentLength", meta.ContentLength, n)
		if meta.ContentLength > 0 {
			errTags := countTags(urlContext)
			if errTags != nil {
				return errTags
			}
			for k, v := range urlContext.tags {
				urlContext.Elements = append(urlContext.Elements, &tagCounter{TagName: k, Count: v})
			}
		}
	}
	log.Trace("End handle url:", urlContext.URL)
	return nil
}

func (resolver *Resolver) resolve(context *HandleContext) error {
	log := resolver.log
	i := 0
	for index, url := range context.urlsToResolve {
		urlContext := context.urlContext[index]
		urlContext.URL = url
		i++
		log.Debug("Send job:", index)
		resolver.jobsChannel <- func() {
			_ = resolver.resolveOneURL(urlContext, resolver.waitResponseTimeoutInMs)
			resolver.log.Trace(urlContext.URL, "answered")
			context.answerChannel <- 1
		}
	}
	log.Debug("Added:", i)
	for {
		select {
		case <-context.answerChannel:
			context.answeredCount++
			log.Trace("context.answeredCount", context.answeredCount, i)
			if context.answeredCount >= len(context.urlsToResolve) {
				log.Trace("end")
				return nil
			}
		}
		log.Trace("Wait")
	}
}

func (resolver *Resolver) getUrls(r *http.Request, context *HandleContext) error {
	n, errRead := r.Body.Read(context.request)
	if errRead != nil && errRead != io.EOF {
		return errors.Wrap(errRead, "Read error. length:"+strconv.FormatInt(int64(n), 10))
	}
	r.Body.Close()
	context.request = context.request[0:n]
	if len(context.request) == 0 {
		return errors.New("Incorect request format. Length is 0")
	}
	if err := json.Unmarshal(context.request, &context.urlsToResolve); err != nil {
		return errors.New("Incorect request format. Unmarshal error")
	}
	resolver.log.Debug("context.urlsToResolve", len(context.urlsToResolve))
	context.urlContext = context.urlContext[0:len(context.urlsToResolve)]
	return nil
}

func (resolver *Resolver) handleResolve(w http.ResponseWriter, r *http.Request) {
	val := atomic.AddInt64(&resolver.countRequests, 1)
	log := resolver.log
	log.Info("Requests:", val)
	<-resolver.limitChannel
	index := 0
	for i, val := range resolver.contextBusy {
		if swapped := atomic.CompareAndSwapInt32(val, 0, 1); swapped {
			index = i
			break
		}
	}

	if resolver.contextes[index] == nil {
		resolver.contextes[index] = NewHandleContext(log)
	}
	context := resolver.contextes[index]
	defer func() {
		context.clear()
		atomic.StoreInt32(resolver.contextBusy[index], 0)
		val := atomic.AddInt64(&resolver.countRequests, -1)
		resolver.limitChannel <- true
		log.Info("Requests:", val)
	}()
	if err := resolver.getUrls(r, context); err != nil {
		log.Debug(err)
		http.Error(w, "Incorect request format.", int(httpBadRequest))
	} else if err = resolver.resolve(context); err != nil {
		log.Debug("resolve error", err)
		http.Error(w, "Incorect request format.", int(httpBadRequest))
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := json.NewEncoder(w).Encode(context.urlContext); err != nil {
			log.Debug("Encode error:", err)
			http.Error(w, "Incorect request format.", int(httpBadRequest))
		}
	}

}

func main() {
	bindAddr := flag.String("bind", "", "address to bind")
	flag.Parse()
	if bindAddr == nil || *bindAddr == "" {
		flag.Usage()
		os.Exit(1)

	}
	logger := NewLogger("log.log", true, false, false)
	limitOfCuncurrentRequests := 2000
	resolver := NewResolver(time.Duration(10*1000*time.Millisecond), limitOfCuncurrentRequests, logger)
	resolver.jobsChannel = make(chan func(), 1000000)
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1) // numCPUs hot threads + one for async tasks.
	var waitGroup sync.WaitGroup
	for i := 0; i < numCPUs; i++ {
		logger.Info("Start thread:", i)
		go func() {
			waitGroup.Add(1)

			for {
				job, ok := <-resolver.jobsChannel
				if !ok {
					logger.Debug("Thread stopped")
					break
				}
				logger.Debug("Job received")
				job()
			}
			waitGroup.Done()
		}()
	}
	http.HandleFunc("/resolve", resolver.handleResolve)
	defer close(resolver.jobsChannel)
	defer waitGroup.Wait()
	if err := http.ListenAndServe(*bindAddr, nil); err != nil {
		log.Println("Error on serve:", err)
		os.Exit(1)
	}

}
