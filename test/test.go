package main

import (
	"github.com/biluping/go-router/router"
	"net/http"
)

func main() {
	router.Get("/hello", func(request *http.Request) interface{} {
		a := 1 - 1
		b := 10 / a
		return b
	})
	router.Start(8080)
}
