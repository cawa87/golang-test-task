package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

const (
	RequestTimeout = time.Second * 60
)

func fetchHandler(c *gin.Context) {

	var request Request

	// decode request
	err := c.BindJSON(&request)
	if err != nil {
		log.Println("Error decoding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	result, err := doFetching(request)
	if err != nil {
		log.Println("Unrecoverable error while fetching the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, result)
}

func doFetching(urls []string) (*Response, error) {

	var (
		result = Response(make([]*ResponseItem, len(urls)))
		wg     sync.WaitGroup
	)

	wg.Add(len(result))

	for i, param := range urls {
		go func(id int, url string) {
			// vars are passed as parameters to create
			// copies from loop ones

			// do the fetching
			res, err := fetchWorker(url)
			if err != nil {
				// error feedback is wanted
				res = &ResponseItem{
					URL: url,
					Meta: Meta{
						Status: http.StatusInternalServerError,
						Error:  err.Error(),
					},
				}
			}

			// arrays are thread-safe until areas
			// do not overlap
			result[id] = res

			wg.Done()
		}(i, param)
	}

	wg.Wait()
	return &result, nil
}

func fetchWorker(url string) (*ResponseItem, error) {

	// do GET with timeout
	client := http.Client{
		Timeout: RequestTimeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	var result = ResponseItem{
		URL: url,
		Meta: Meta{
			Status: resp.StatusCode,
		},
	}

	// fill in content type if code is 2xx
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		err = resp.Body.Close()
		return &result, err
	}

	result.Meta.ContentType = resp.Header.Get("Content-Type")

	// read body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result.Meta.ContentLength = len(body)

	// abort if content is empty or not html
	if result.Meta.ContentLength == 0 ||
		!strings.HasPrefix(result.Meta.ContentType, "text/html") {

		return &result, nil
	}

	// count && fill-in tags
	tags, err := countTags(body)
	if err != nil {
		return nil, err
	}
	result.Elements = tags

	return &result, nil
}

// count all html-tags in input document
//  <p>lol</p> is one p element
func countTags(body []byte) ([]Element, error) {

	var (
		counts = map[string]int{}
		reader = bytes.NewBuffer(body)
		z      = html.NewTokenizer(reader)
	)

	for {
		switch z.Next() {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				// this is the return-point
				return encodeTags(counts), nil
			}

			// any other err is unexpected
			log.Printf("html token err: %v", z.Err())
			return nil, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, _ := z.TagName()
			counts[string(tagName)] += 1
		}
	}

}

// encode counts map to a response array
func encodeTags(counts map[string]int) []Element {
	var (
		result = make([]Element, len(counts))
		iter   int
	)

	// decode it to result
	for name, count := range counts {
		result[iter] = Element{
			TagName: name,
			Count:   count,
		}
		iter++
	}

	return result
}
