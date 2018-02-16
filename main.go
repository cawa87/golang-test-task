package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"bytes"

	"golang.org/x/net/html"
)

func main() {
	http.HandleFunc("/", parseUrls)
	http.ListenAndServe(":8000", nil)
	log.Println("server is running on :8000")
}

type parseResult struct {
	URL      string               `json:"url"`
	Meta     parseResultMeta      `json:"meta"`
	Elements []parseResultElement `json:"elements,omitempty"`
}

type parseResultMeta struct {
	Status        int    `json:"status"`
	ContentLength int    `json:"content-length,omitempty"`
	ContentType   string `json:"content-type,omitempty"`
}

type parseResultElement struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

func parseUrls(w http.ResponseWriter, r *http.Request) {
	urls := parseRequest(r)

	resultChannel := make(chan parseResult)
	defer close(resultChannel)

	for _, url := range urls {
		go parseURL(url, resultChannel)
	}

	results := make([]parseResult, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		result := <-resultChannel
		results = append(results, result)
	}

	body, _ := json.Marshal(results)

	io.WriteString(w, string(body))
}

func parseRequest(r *http.Request) []string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	var urls []string
	_ = json.Unmarshal(body, &urls)
	return urls
}

func parseURL(url string, resultChannel chan<- parseResult) {
	resp, _ := http.Get(url)
	if resp == nil || resp.Body == nil {
		resultChannel <- parseResult{
			URL: url,
		}
		return
	}
	if resp.Body == nil {
		resultChannel <- parseResult{
			URL: url,
			Meta: parseResultMeta{
				Status: resp.StatusCode,
			},
		}
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	tagsMap := map[string]int{}
	tokenizer := html.NewTokenizer(bytes.NewReader(body))

	active := true
	for active {
		tokenType := tokenizer.Next()
		tag, _ := tokenizer.TagName()

		switch tokenType {
		case html.ErrorToken:
			active = false
		case html.StartTagToken, html.SelfClosingTagToken:
			tagsMap[string(tag)] += 1
		}
	}

	elements := make([]parseResultElement, 0, len(tagsMap))
	for name, count := range tagsMap {
		elements = append(elements, parseResultElement{
			TagName: name,
			Count:   count,
		})
	}

	resultChannel <- parseResult{
		URL: url,
		Meta: parseResultMeta{
			Status:        resp.StatusCode,
			ContentLength: len(body), // because resp.ContentLength is -1
			ContentType:   resp.Header.Get("Content-Type"),
		},
		Elements: elements,
	}
}
