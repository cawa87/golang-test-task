package main

import "er"
import "mime"
import "net/http"
import "sync/atomic"

// Any task's always being processed exclusively
type FetchAndScanTask struct {
	context *FetchAndScanContext
	data    *FetchAndScanData
}

type FetchAndScanContext struct {
	result chan<- *FetchAndScanData
	canceled int32                  // Shared
}

type FetchAndScanData struct {
	Url     string                     `json:"url"`
	Meta    FetchAndScanDataMeta       `json:"meta,omitempty"`
	Elemets []FetchAndScanDataElemets  `json:"elemets,omitempty"`

	e error
	response *http.Response
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

func (t FetchAndScanTask) Done(e error) {
	if t.data.e == nil { t.data.e = e }
	t.context.result <- t.data
}

func (c *FetchAndScanContext) Cancel() {
	atomic.SwapInt32(&c.canceled, 1)
}

func (c *FetchAndScanContext) IsCanceled() bool {
	return atomic.LoadInt32(&c.canceled) != 0
}

func (d *FetchAndScanData) Release() {
	if d.response != nil {
		d.response.Body.Close()
		d.response = nil
	}
}

func startFetchWorkers(workers int, pipe <-chan FetchAndScanTask)  {
	for i := 0; i < workers; i++ { go fetchWorker(pipe) }
}

func headerGetMediaType(header http.Header, key string) string {
	if a := header.Get(key); a != "" {
		t, _, e := mime.ParseMediaType(a)
		if e == nil { return t }}
	return ""
}

func fetchWorker(pipe <-chan FetchAndScanTask) {
	for task := range pipe {
		if task.context.IsCanceled() { task.Done(nil); continue }

		response, e := http.Get(task.data.Url)
		if e != nil { task.Done(er.Er(e, "http.Get")); continue }

		task.data.response = response
		task.data.Meta.Status = response.StatusCode
		if ct := headerGetMediaType(response.Header, "content-type"); ct != "" {
			task.data.Meta.ContentType = &ct }
		if cl := response.ContentLength; cl >= 0 {
			task.data.Meta.ContentLength = &cl }

		task.Done(nil)
	}
}

func fetchAndScan(
	urls []string,
	fetchPipe chan<- FetchAndScanTask,
) (result []*FetchAndScanData, E error) {
	result = []*FetchAndScanData{}
	if len(urls) < 1 { return }

	output := make(chan *FetchAndScanData)
	defer close(output)

	context := FetchAndScanContext{output, 0}
	go func() {
		for _, url := range urls {
			data := FetchAndScanData{}
			data.Url = url
			fetchPipe <- FetchAndScanTask{&context, &data}
		}
	}()

	processed := 0
	for data := range output {
		data.Release()

		if data.e != nil {
			E = data.e
			context.Cancel()
		} else {
			result = append(result, data)
		}

		processed += 1
		if processed >= len(urls) { break }
	}
	return
}
