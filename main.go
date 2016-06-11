package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

const mdTempl = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Marked in the browser</title>
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.6/css/bootstrap.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.4.0/styles/github.min.css">
  <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.0.0/jquery.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/marked/0.3.5/marked.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.4.0/highlight.min.js"></script>
  <script>
$(document).ready(function(){
  hljs.initHighlightingOnLoad();
  marked.setOptions({ langPrefix: '' });
  var target = $("#markdown_content");
  $.ajax({
    url: "{{.RequestURI}}?raw=1",
  }).done(function(data){
    target.append(marked(data));
  }).fail(function(data){
    target.append("This content failed to load.");
  });
});
  </script>
</head>
<body>
  <!-- Content -->
  <div class="container">
    <div id="markdown_content"> </div>
  </div>
</body>
</html>
`

type SuffixMux struct {
	m          map[string]http.Handler
	defHandler http.Handler
}

func NewSuffixMux() *SuffixMux {
	return &SuffixMux{m: make(map[string]http.Handler)}
}

func (mux *SuffixMux) handler(r *http.Request) http.Handler {
	for s, h := range mux.m {
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
	t := template.Must(template.New("letter").Parse(mdTempl))
	t.Execute(w, r)
}

func main() {
	mux := NewSuffixMux()
	mux.Handle(".md", http.HandlerFunc(MarkDownHandler))
	mux.DefaultHandler(http.FileServer(http.FileSystem(http.Dir("."))))
	log.Println("Listening at 0.0.0.0:8888")
	http.ListenAndServe("0.0.0.0:8888", mux)
}
