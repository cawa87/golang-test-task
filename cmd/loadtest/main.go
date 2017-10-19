package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
)

type jobResult struct {
	resp *http.Response
	err  error
}

func job(jobs chan []string, results chan<- jobResult) {
	for urls := range jobs {
		buffer := new(bytes.Buffer)

		json.NewEncoder(buffer).Encode(urls)

		resp, err := http.Post("http://localhost:9000/v1/", "application/json", buffer)

		results <- jobResult{
			resp: resp,
			err:  err,
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("path to csv required")
		os.Exit(1)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		os.Exit(0)
	}()

	jobs := make(chan []string, 100)
	results := make(chan jobResult, 100)

	for w := 0; w < 100; w++ {
		go job(jobs, results)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	csvReader := csv.NewReader(bufio.NewReader(file))

	urls := []string{}
	jobsCount := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		urls = append(urls, record[0])

		if len(urls) >= 50 {
			jobsCount++
			jobs <- urls
			urls = []string{}
		}

		if jobsCount >= 30 {
			break
		}
	}
	file.Close()

	for i := 1; i <= jobsCount; i++ {
		r := <-results
		if r.err != nil {
			fmt.Println(r.err)
		} else {
			d, _ := httputil.DumpResponse(r.resp, true)
			fmt.Println(string(d))
		}
		if r.resp != nil {
			r.resp.Body.Close()
		}
	}
}
