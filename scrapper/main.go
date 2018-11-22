package main

import (
	"flag"
	"github.com/souz9/golang-test-task/scrapper/hchecker"
	"github.com/souz9/golang-test-task/scrapper/server"
	"log"
	"time"
)

func main() {
	var (
		listen        = flag.String("l", ":8080", "Listen requests on")
		sitesPath     = flag.String("s", "scrapper/sites.txt", "List of sites to watch")
		watchInterval = flag.Duration("w", time.Minute, "Watch interval")
	)
	flag.Parse()

	sites, err := readLines(*sitesPath)
	if err != nil {
		log.Fatalf("read sites: %v", err)
	}

	hc := hchecker.New()
	hc.Watch(sites, *watchInterval)

	s := server.Server{HChecker: hc}
	log.Fatal(s.Listen(*listen))
}
