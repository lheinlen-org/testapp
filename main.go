package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	page = `<html>
<head>
  <title>Hello World!</title>
</head>
<body>
  <p>
	{{.Message}}
  </p>
</body>
</html>`
)

var (
	tmpl = template.Must(template.New("page").Parse(page))
)

type templateHeader struct {
	Name  string
	Value string
}

type templateData struct {
	Headers []templateHeader
	Message string
}

func main() {
	port := os.Getenv("HWW_PORT")
	if port == "" {
		port = "8080"
	}

	msg := os.Getenv("HWW_MESSAGE")
	if msg == "" {
		msg = "Hello World!"
	}

	log.Printf("Starting server, port %s, message: %s", port, msg)
	http.HandleFunc("/", getHandler(msg))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server stopped, error: %s", err)
	}
}

func getHandler(msg string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a request, source: %s", r.RemoteAddr)

		data := templateData{Message: msg}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template failed: %s", err.Error())
			http.Error(w, err.Error(), 500)
		}
	}
}
