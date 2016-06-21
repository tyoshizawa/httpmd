# httpmd: Static http server with Markdown support
## What's this
httpmd is a Golang simple static web server based on net/http.FileServer with
* Client side Markdown rendering using [marked.js](https://github.com/chjj/marked) when accessing to Markdown file with render=1 query string.
* Client side source code highliting as a Markdown code block when accessing to source code file with render=1 query string.
* Modified index page with render=1 query string for Markdown and source code files link URL.

## How to Install
```shell
go get github.com/tyoshizawa/httpmd
```

## How to Use
```
httpmd [-h host] [-p port] [-d dir]
  -h host : Host ip address. Default is 0.0.0.0.
  -p port : Port number for the server. Default is 8888.
  -d dir  : Document root directory. Default is ".".
```

