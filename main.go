package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"strings"

	"github.com/nrvru/golang-test-task/models"
	"golang.org/x/net/html"
)

func root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			io.WriteString(w, fmt.Sprintf("Media type '%s' not supported\n", contentType))
			return
		}

		handleUrls(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, fmt.Sprintf("Method '%s' not allowed\n", r.Method))
	}
}

func handleUrls(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var urls []string

	if err := d.Decode(&urls); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	result := models.Result{}
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			var err error
			urlData, err := crawl(url)
			if err != nil {
				log.Println(err.Error())
			}

			if urlData != nil {
				result.Add(urlData)
			}

			wg.Done()
		}(url)
	}
	wg.Wait()

	data, err := result.JSON()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("Internal sever error\n"))
		return
	}

	w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func crawl(url string) (*models.UrlData, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Can not get url '%s': %s\n", url, err.Error())
	}
	defer res.Body.Close()

	data := &models.UrlData{
		Url: url,
		Meta: models.Meta{
			Status: res.StatusCode,
		},
	}

	if res.StatusCode != http.StatusOK {
		return data, nil
	}

	counter := &models.LengthCounter{}
	dataSrc := io.TeeReader(res.Body, counter)

	contentType := res.Header.Get("Content-Type")
	data.Meta.ContentType = contentType

	if !strings.Contains(contentType, "text/html") {
		io.Copy(ioutil.Discard, dataSrc)
		data.Meta.ContentLength = counter.Total
		return data, nil
	}

	doc := html.NewTokenizer(dataSrc)
	elements := map[string]int{}

	for {
		tt := doc.Next()
		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := doc.Token()
			if c, ok := elements[t.Data]; ok {
				elements[t.Data] = c + 1
			} else {
				elements[t.Data] = 1
			}
		}
	}

	data.Meta.ContentLength = counter.Total

	for tag, count := range elements {
		data.Elements = append(data.Elements, models.Element{
			TagName: tag,
			Count:   count,
		})
	}

	return data, nil
}

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}
