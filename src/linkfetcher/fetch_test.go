package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

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

	res, err := doFetching([]string{
		"https://linux.org.ru",
	})

	assert.Nil(t, err)
	assert.Equal(t, res, excepted)

}
