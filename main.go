package main

import (
	"log"
	"net/http"
)

func main () {
	router := http.NewServeMux()

	router.HandleFunc("/{$}",Home)
	router.HandleFunc("/projects",ListProjects)
	router.HandleFunc("/projects/{id}",ShowProjectsById)
	router.HandleFunc("/",NotFound)

	fileServer := http.FileServer(http.Dir("./assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

