package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Println("Starting up")

	r := gin.Default()

	r.POST("/fetch", fetchHandler)

	log.Fatal(r.Run())
}
