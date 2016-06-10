package main

import (
	"net/http"
	"log"
	"strings"
)

type SuffixMux struct{
	m map[string]http.Handler
	defHandler http.Handler
}

func NewSuffixMux() *SuffixMux {
	return &SuffixMux{m: make(map[string]http.Handler)}
}

func (mux *SuffixMux) handler(r *http.Request) http.Handler {
	for s, h := range(mux.m) {
		if strings.HasSuffix(r.RequestURI, s) {
			return h
		}
	}
	return mux.defHandler
}

func (mux *SuffixMux) Handle(suffix string, h http.Handler) {
	mux.m[suffix] = h
}

func (mux *SuffixMux) DefaultHandler(h http.Handler) {
	mux.defHandler = h
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
	mux.Handle(".md", http.HandlerFunc(MarkDownHandler))
	mux.DefaultHandler(http.FileServer(http.FileSystem(http.Dir("."))))
	log.Println("Listening at 0.0.0.0:8888")
	http.ListenAndServe("0.0.0.0:8888", mux)
}
