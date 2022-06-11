package main

import (
	"github.com/biluping/go-router/router"
)

type User struct {
	Name string
	Age  string
}

func main() {
	router.Post("/hello", "aaa")
	router.Start(8080)
}
