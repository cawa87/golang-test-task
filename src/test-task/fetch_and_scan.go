package main

import "bytes"
import "er"
import "golang.org/x/net/html"
import "http_plus"
import "io"
import "io/ioutil"
import "net/http"
import "sort"
import "sync/atomic"
import "time"

type FetchAndScan struct {
	fetchPipe chan *FetchAndScanTask
	scanPipe  chan *FetchAndScanTask
}

// Any task's always processed exclusively
type FetchAndScanTask struct {
	context *FetchAndScanContext
	data    FetchAndScanData
	e       error
	body    []byte
}

type FetchAndScanContext struct {
	done     chan<- *FetchAndScanTask
	canceled int32 // Shared
}

type FetchAndScanData struct {
	Url     string                    `json:"url"`
	Meta    FetchAndScanDataMeta      `json:"meta,omitempty"`
	Elements []FetchAndScanDataElement `json:"elemets,omitempty"`
}

type FetchAndScanDataMeta struct {
	Status        int     `json:"status"`
	ContentType   *string `json:"content-type,omitempty"`
	ContentLength *int64  `json:"content-length,omitempty"`
}

type FetchAndScanDataElement struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

func NewFetchAndScan(
	fetchWorkers, scanWorkers, taskBufferSize int) *FetchAndScan {
	fetchPipe := make(chan *FetchAndScanTask, taskBufferSize)
	scanPipe  := make(chan *FetchAndScanTask, taskBufferSize)
	for i := 0; i < fetchWorkers; i++ { go fetchWorker(fetchPipe, scanPipe) }
	for i := 0; i < scanWorkers;  i++ { go scanWorker(scanPipe) }
	return &FetchAndScan{fetchPipe, scanPipe}
}

func NewFetchAndScanTask(
	context *FetchAndScanContext, url string) *FetchAndScanTask {
	task := &FetchAndScanTask{}
		task.context = context
		task.data.Url = url
	return task
}

func (fs *FetchAndScan) Do(urls []string) (result []FetchAndScanData, E error) {
	result = []FetchAndScanData{}
	if len(urls) < 1 { return }

	done    := make(chan *FetchAndScanTask)
	context := FetchAndScanContext{done, 0}
	go func() {
		for _, url := range urls {
			fs.fetchPipe <- NewFetchAndScanTask(&context, url) }
	}()

	processed := 0
	for task := range done {
		if task.e != nil {
			E = task.e
			context.Cancel()
		} else {
			result = append(result, task.data)
		}

		processed += 1
		if processed >= len(urls) { break }
	}
	return
}

func (t *FetchAndScanTask) Done(e error) {
	if t.e == nil { t.e = e }
	t.context.done <- t
}

func (c *FetchAndScanContext) Cancel() {
	atomic.SwapInt32(&c.canceled, 1)
}

func (c *FetchAndScanContext) IsCanceled() bool {
	return atomic.LoadInt32(&c.canceled) != 0
}

func fetchWorker(
	pipe <-chan *FetchAndScanTask, scanPipe chan<- *FetchAndScanTask,
){
	for task := range pipe {
		if task.context.IsCanceled() { task.Done(nil); continue }

		// Request Url
		httpClient := http.Client{Timeout: 10 * time.Second}
		req, e := http.NewRequest("GET", task.data.Url, nil)
		if e != nil {
			task.Done(er.Er(e, "http.NewRequest", "url", task.data.Url))
			continue
		}
		req.Header.Set("Accept-Encoding", "") // We need original content, not compressed

		res, e := httpClient.Do(req)
		if e != nil {
			task.Done(er.Er(e, "http.Client.Do", "url", task.data.Url))
			continue
		}
		defer res.Body.Close()

		// Fill Status, ContentType, ContentLength
		sc := res.StatusCode
			task.data.Meta.Status = sc
		ct := http_plus.HeaderGetMediaType(res.Header, "content-type")
			if ct != "" { task.data.Meta.ContentType = &ct }
		cl := http_plus.HeaderGetContentLength(res.Header, res.ContentLength)
			if cl >= 0 { task.data.Meta.ContentLength = &cl }

		// Read body if text/html and forward for scanning
		if 200 <= sc && sc < 300 && ct == "text/html" {
			task.body, e = ioutil.ReadAll(res.Body)
			if e != nil {
				task.Done(er.Er(e, "Fail to read body", "url", task.data.Url))
				continue
			}
			scanPipe <- task
			continue
		}

		task.Done(nil)
	}
}

func scanWorker(pipe <-chan *FetchAndScanTask) {
	worker: for task := range pipe {
		if task.context.IsCanceled() { task.Done(nil); continue }

		// Update ContentLength in case of the original header missed or was invalid
		cl := int64(len(task.body))
		task.data.Meta.ContentLength = &cl

		// Count HTML tags
		tags := map[string]int{}
		z := html.NewTokenizer(bytes.NewReader(task.body))
		for {
			t, e := z.Next(), z.Err()
			if e == io.EOF { break }
			if e != nil    { task.Done(er.Er(e, "Fail to parse html")); continue worker }
			switch t { case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
				tag, _ := z.TagName()
				tags[string(tag)]++
			}
		}

		// Tags to Elements
		elements := []FetchAndScanDataElement{}
		for tag, count := range tags {
			elements = append(elements,
				FetchAndScanDataElement{tag, count})
		}
		sort.Slice(elements, func(a, b int) bool {
			return elements[a].TagName < elements[b].TagName })
		task.data.Elements = elements

		task.Done(nil)
	}
}
