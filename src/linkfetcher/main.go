package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Println("Starting up")

	// create fetcher back-end
	fetcher, err := newFetcher()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	// register your routes in there
	r.POST("/fetch", fetcher.handle)

	// serve forever
	// @TODO: a good microservice will listen for
	//  signals and stop everything gracefully
	log.Fatal(r.Run())
}
