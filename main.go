package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"se.com/data-processing/devdevman/processing"
)

const basePath = "/api"

func main() {
    r := mux.NewRouter()

	processing.SetupRoutes(basePath, r)
    http.Handle("/", r)
	log.Printf("Visit http://localhost:%d/ for giggles", 5000)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
