package main

import (
	"log"
	"net/http"
)

func main() {

	router := NewRouter()

	if err := initDocumets("docs"); err != false {
	}

	log.Fatal(http.ListenAndServe(":4730", router))
}
