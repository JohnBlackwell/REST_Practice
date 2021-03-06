package main

import (
	"github.com/go-martini/martini"
)

func main() {
	//set up server
	m := martini.Classic()
	m.Get("/", func() string {
		return "Running on localhost"
	})
	m.Post("/process", ProcessRequest)
	m.Run()
}
