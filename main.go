package main

import (
	"net/http"
	"log"
	"strings"
)

type SuffixMux struct{
	defHandler http.Handler
}

func NewSuffixMux() *SuffixMux {
	return &SuffixMux{}
}

func (mux *SuffixMux) handler(r *http.Request) http.Handler {
	if strings.HasSuffix(r.RequestURI, ".md") {
		return http.HandlerFunc(MarkDownHandler)
	} else {
		return mux.defHandler
	}
}

func (mux *SuffixMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h := mux.handler(r)
	h.ServeHTTP(w, r)
}

func MarkDownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {
	mux := NewSuffixMux()
	mux.defHandler = http.FileServer(http.FileSystem(http.Dir(".")))
	log.Println("Listening at 0.0.0.0:8888")
	http.ListenAndServe("0.0.0.0:8888", mux)
}
