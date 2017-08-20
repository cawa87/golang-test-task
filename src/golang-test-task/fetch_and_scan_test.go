package main

import "encoding/json"
import "errors"
import "fmt"
import "github.com/stretchr/testify/assert"
import "github.com/xeipuuv/gojsonschema"
import "gopkg.in/jarcoal/httpmock.v1"
import "net/http"
import "testing"

const FETCH_AND_SCAN_DATA_SCHEMA = `
{
	"type": "array",
	"items": {
		"type": "object",
		"required": ["url", "meta"],
		"properties": {
			"url": {
				"type": "string",
				"format": "uri",
				"description": "uri from input list"
			},
			"meta": {
				"type": "object",
				"required": ["status"],
				"properties": {
					"status": {
							"type": "integer",
							"description": "Response status of this uri"
					},
					"content-type": {
							 "type": "string",
							 "description": "In case of 2XX response status, value of mime-type part of Content-Type header (if exists)"
					},
					"content-length": {
								"type": "integer",
								"minimum": 0,
								"description": "In case of 2XX response status, length of response body (be careful, response could be chunked)."
					}
				}
			},
			"elemets": {
				"type": "array",
				"description": "In case of 2XX response status, \"text\/html\" content type and non-zero content length, list of HTML-tags, occured.",
				"items": {
					"type": "object",
					"required": ["tag-name", "count"],
					"properties": {
						"tag-name": {"type": "string"},
						"count": {
							"type": "integer",
							"minimum": 1,
							"description": "Number of times, the current tag occures in response"
						}
					}
				}
			}
		}
	}
}
`
var fetchAndScanDataSchema *gojsonschema.Schema

var _ = func() (_ struct{}) {
	var e error
  a := gojsonschema.NewStringLoader(FETCH_AND_SCAN_DATA_SCHEMA)
	fetchAndScanDataSchema, e = gojsonschema.NewSchema(a)
	if e != nil { panic(e) }
	return
}()

func fetchAndScanValidateSchema(data string) error {
	a := gojsonschema.NewStringLoader(data)
	b, e := fetchAndScanDataSchema.Validate(a)
	if e != nil { return e }
	if !b.Valid() {
		desc := ""
		for _, e := range b.Errors() {
			desc += fmt.Sprintf("[%s: %s]", e.Field(), e.Description()) }
		msg := fmt.Sprintf("Invalid schema desc=%s data=%s", desc, data)
		return errors.New(msg)
	}
	return nil
}

func makeFuncAndScanContext(
	fetchWorkers, scanWorkers, taskBufferSize int,
)(
	fetchAndScan *FetchAndScan,
	context      *FetchAndScanContext,
	done         chan *FetchAndScanTask,
){
	fetchAndScan = NewFetchAndScan(fetchWorkers, scanWorkers, taskBufferSize)
	done         = make(chan *FetchAndScanTask)
	context      = &FetchAndScanContext{done, 0}
	return
}

func TestFetchWorkerShouldAdvance(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(2, 0, 0)
	for i := 0; i < 3; i++ {
		fs.fetchPipe <- NewFetchAndScanTask(context, "")
		assert.Error(t, (<-done).e)
	}
}

func TestFetchWorkerShouldAdvanceCanceled(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(2, 0, 0)
	context.Cancel()
	for i := 0; i < 3; i++ {
		fs.fetchPipe <- NewFetchAndScanTask(context, "")
		assert.NoError(t, (<-done).e)
	}
}

func TestScanWorkerShouldAdvance(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(0, 2, 0)
	for i := 0; i < 3; i++ {
		fs.scanPipe <- NewFetchAndScanTask(context, "")
		assert.NoError(t, (<-done).e)
	}
}

func TestScanWorkerShouldAdvanceCanceled(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(0, 2, 0)
	context.Cancel()
	for i := 0; i < 3; i++ {
		fs.scanPipe <- NewFetchAndScanTask(context, "")
		assert.NoError(t, (<-done).e)
	}
}

func testFetchAndScanResult(
	t *testing.T, fetchAndScan *FetchAndScan, urls []string, expected string) {
	data, e := fetchAndScan.Do(urls)
	assert.NoError(t, e)
	bytes, e := json.Marshal(&data)
	assert.NoError(t, e)
	assert.NoError(t, fetchAndScanValidateSchema(string(bytes)))
	assert.Equal(t, expected, string(bytes))
}

func TestFetchAndScanEmptyUrls(t *testing.T) {
	testFetchAndScanResult(t, NewFetchAndScan(0, 0, 0),
		[]string{}, `[]`)
}

func TestFetchAndScanUnknownUrl(t *testing.T) {
	_, e := NewFetchAndScan(1, 0, 0).Do([]string{"http://unknown"})
	assert.Error(t, e)
}

func TestFetchAndScanUrlWithNoContentLengthButData(t *testing.T) {
	httpmock.Activate(); defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://one",
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(200, "<a></a>")
			res.Header.Set("Content-Length", "")
			return res, nil
		})
	testFetchAndScanResult(t, NewFetchAndScan(1, 1, 0),
		[]string{"https://one"},
		`[{"url":"https://one","meta":{"status":200}}]`,
	)
}

func TestFetchAndScanUrlWithNoContentLengthButHtmlData(t *testing.T) {
	httpmock.Activate(); defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://one",
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(200, "<a></a>")
			res.Header.Set("Content-Type", "text/html")
			res.Header.Set("Content-Length", "")
			return res, nil
		})
	testFetchAndScanResult(t, NewFetchAndScan(1, 1, 0),
		[]string{"https://one"},
		`[{"url":"https://one","meta":{"status":200,"content-type":"text/html","content-length":7},"elemets":[{"tag-name":"a","count":2}]}]`,
	)
}

func TestFetchAndScanUrlWithHtmlButNot2XX(t *testing.T) {
	httpmock.Activate(); defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://one",
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(123, "<a></a>")
			res.Header.Set("Content-Type", "text/html")
			res.Header.Set("Content-Length", "7")
			return res, nil
		})
	testFetchAndScanResult(t, NewFetchAndScan(1, 1, 0),
		[]string{"https://one"},
		`[{"url":"https://one","meta":{"status":123,"content-type":"text/html","content-length":7}}]`,
	)
}

func TestFetchAndScanUrlWithEmptyHtml(t *testing.T) {
	httpmock.Activate(); defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://one",
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(200, "")
			res.Header.Set("Content-Type", "text/html")
			return res, nil
		})
	testFetchAndScanResult(t, NewFetchAndScan(1, 1, 0),
		[]string{"https://one"},
		`[{"url":"https://one","meta":{"status":200,"content-type":"text/html","content-length":0}}]`,
	)
}

func TestFetchAndScanUrlWithHtml(t *testing.T) {
	httpmock.Activate(); defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://one",
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(200,
				"<a></a><a/>b<!--c--><d>")
			res.Header.Set("Content-Type", "text/html")
			return res, nil
		})
	testFetchAndScanResult(t, NewFetchAndScan(1, 1, 0),
		[]string{"https://one"},
		`[{"url":"https://one","meta":{"status":200,"content-type":"text/html","content-length":23},"elemets":[{"tag-name":"a","count":3},{"tag-name":"d","count":1}]}]`,
	)
}
