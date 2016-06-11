package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
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
  marked.setOptions({ langPrefix: '' });
  var target = $("#markdown_content");
  $.ajax({
    url: "{{.URI}}?raw=1",
    dataType: "text",
  }).done(function(data){
    {{if .LANG}}data = "## {{.URI}}\n` + "```" + `{{.LANG}}\n" + data + "\n` + "```" + `";{{end}}
    target.append(marked(data));
    $('#markdown_content pre code').each(function(i, block) {
      hljs.highlightBlock(block);
    });
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

// codes is a map of key: suffix string, value: lang
var codes = map[string]string {
	".c": "c",
	".cpp": "cpp",
	".css": "css",
	".diff": "diff",
	".go": "go",
	".java": "java",
	".js": "javascript",
	".json": "json",
	".pl": "perl",
	".php": "php",
	".py": "python",
	".rb": "ruby",
	".sh": "shell",
	".sql": "sql",
}

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

func CodeMarkDownHandler(lang string) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("markdown").Parse(mdTempl))
		uri := template.JSEscapeString(r.RequestURI)
		t.Execute(w, struct{ URI, LANG string }{uri, lang})
	}
	return http.HandlerFunc(handler)
}

func main() {
	mux := NewSuffixMux()
	mux.Handle(".md", CodeMarkDownHandler(""))
	for sfx, lang := range codes {
		mux.Handle(sfx, CodeMarkDownHandler(lang))
	}
	mux.DefaultHandler(http.FileServer(http.FileSystem(http.Dir("."))))
	log.Println("Listening at 0.0.0.0:8888")
	http.ListenAndServe("0.0.0.0:8888", mux)
}
