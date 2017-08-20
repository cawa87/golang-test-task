package main

import "er"
import "net/http"
import "sync/atomic"

type FetchAndScan struct {
	fetchPipe chan *FetchAndScanTask
	scanPipe  chan *FetchAndScanTask
}

// Any task's always processed exclusively
type FetchAndScanTask struct {
	context *FetchAndScanContext
	data    FetchAndScanData
	e       error
}

type FetchAndScanContext struct {
	done     chan<- *FetchAndScanTask
	canceled int32 // Shared
}

type FetchAndScanData struct {
	Url     string                     `json:"url"`
	Meta    FetchAndScanDataMeta       `json:"meta,omitempty"`
	Elemets []FetchAndScanDataElemets  `json:"elemets,omitempty"`
}

type FetchAndScanDataMeta struct {
	Status        int     `json:"status"`
	ContentType   *string `json:"content-type,omitempty"`
	ContentLength *int64  `json:"content-length,omitempty"`
}

type FetchAndScanDataElemets struct {
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

		response, e := http.Get(task.data.Url)
		if e != nil { task.Done(er.Er(e, "http.Get")); continue }
// TODO: response.Body.close()

		data := &task.data
		data.Meta.Status = response.StatusCode
		if ct := headerGetMediaType(response.Header, "content-type"); ct != "" {
			data.Meta.ContentType = &ct }
		if cl := response.ContentLength; cl >= 0 {
			data.Meta.ContentLength = &cl }

		scanPipe <- task
	}
}

func scanWorker(pipe <-chan *FetchAndScanTask) {
	for task := range pipe {
		if task.context.IsCanceled() { task.Done(nil); continue }

		task.Done(nil)
	}
}
