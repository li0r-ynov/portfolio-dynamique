package main

import (
	"log"
	"net/http"
)

func main () {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", Home)
	router.HandleFunc("GET /projects", ListProjects)
	router.HandleFunc("GET /projects/{id}", ShowProjectsById)
	router.HandleFunc("GET /", NotFound)

	fileServer := http.FileServer(http.Dir("./assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

