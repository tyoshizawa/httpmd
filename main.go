package main

import (
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FileSystem(http.Dir("."))))
	log.Println("Listening at 0.0.0.0:8888")
	http.ListenAndServe("0.0.0.0:8888", mux)
}
