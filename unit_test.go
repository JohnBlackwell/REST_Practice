package main

import (
	"fmt"
	"testing"

	"github.com/go-martini/martini"
)

func startServer() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Running on localhost"
	})
	m.Post("/process", ProcessRequest)
	m.Run()
}

func TestDBConnection(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Cannot connect to db %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("Cannot ping db %s", err.Error())
	}
}

/*func TestFileUpload(t *testing.T){

}*/

func TestPasswordChecker(t *testing.T){
	url := "http://localhost:3000/process"
	payload := strings.NewReader("password=hello_world")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		t.Fatalf("%s: ", err.Error())
	}
	req.Header.Add("content-type", "application/json")
}