package main

import (
	"github.com/biluping/go-router/router"
	"net/http"
)

func main() {
	router.Get("/hello", func(request *http.Request) interface{} {
		return "helo"
	})
	router.Start(8080)
}
