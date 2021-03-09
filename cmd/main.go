package main

import (
	"bufio"
	"context"

	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkucherenkov/golang-test-task/controller"
	"github.com/kkucherenkov/golang-test-task/repository"
	"github.com/kkucherenkov/golang-test-task/scrapper"
	"github.com/kkucherenkov/golang-test-task/transport"
)

func main() {

	file, err := os.Open("../scrapper/sites.txt")
	if err != nil {
		log.Fatal(err)
	}

	var hosts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	file.Close()

	repo := repository.New(hosts)
	srv := scrapper.NewHTTPService(repo)

	ctx, cnl := context.WithCancel(context.Background())

	defer cnl()

	// every 30 seconds start polling
	tick := time.NewTicker(time.Second * 30)
	done := make(chan bool)
	go scheduler(ctx, srv, tick, done)

	errs := make(chan error, 2)
	ssc := controller.New(repo)

	fmt.Println("start HTTP server")

	handler, error := transport.MakeHandlerJSONRPC(ssc)
	if error != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		done <- true
	}()

	errs <- fmt.Errorf("server terminated")
}

func scheduler(ctx context.Context, srv scrapper.Service, tick *time.Ticker, done chan bool) {
	task(ctx, srv, time.Now())
	for {
		select {
		case t := <-tick.C:
			task(ctx, srv, t)
		case <-done:
			return
		}
	}
}

func task(ctx context.Context, srv scrapper.Service, t time.Time) {
	fmt.Println("scrapping started at", t)
	_ = srv.StartScrapWithContext(ctx)
	fmt.Println("scrapping finished", time.Now().Unix()-t.Unix())
}
