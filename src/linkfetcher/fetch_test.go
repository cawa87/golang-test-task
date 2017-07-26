package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

// Sample test that uses HTTP mock to fake HTTP request
func TestFetchSampleHTML(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var (
		sampleURL  = "https://linux.org.ru"
		sampleHTML = []byte(`<p>hellow test<br/><br/></p>`)
		excepted   = &Response{
			&ResponseItem{
				URL: sampleURL,
				Meta: Meta{
					Status:        200,
					ContentType:   "text/html",
					ContentLength: len(sampleHTML),
				},
				Elements: []Element{
					{
						TagName: "p",
						Count:   1,
					}, {
						TagName: "br",
						Count:   2,
					},
				},
			},
		}
	)

	// create an excepted HTML response
	response := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(sampleHTML)),
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
	}

	httpmock.RegisterResponder("GET", sampleURL,
		httpmock.ResponderFromResponse(response))

	fs, err := newFetcher()
	assert.Nil(t, err)

	res, err := fs.do([]string{
		"https://linux.org.ru",
	})

	assert.Nil(t, err)
	assert.Equal(t, res, excepted)

}
