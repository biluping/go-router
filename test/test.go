package main

import (
	"github.com/biluping/go-router/router"
	"net/http"
)

func main() {
	router.Get("/hello", func(request *http.Request) interface{} {
		return "hello"
	})
	router.Start(8080)
}
