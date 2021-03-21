package webapp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type indexPage []byte

func (ip indexPage) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(ip)
	if err != nil {
		log.Println(err)
	}
}

type indexPageParams struct {
	Host string
	Port uint16
}

//go:embed index.gohtml
var indexTemplate embed.FS

func newIndexPage(port uint16) indexPage {
	indexPageFs, err := template.ParseFS(indexTemplate, "*.gohtml")
	if err != nil {
		log.Panicf("Can't parse index page; %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	hp := &indexPageParams{
		Host: hostname,
		Port: port,
	}

	buff := &bytes.Buffer{}
	if err = indexPageFs.Execute(buff, hp); err != nil {
		log.Panicf("Can't parse indexPage page; %v", err)
	}

	fmt.Printf("In your web browser, go to http://%s:%d\n", hostname, port)

	return buff.Bytes()
}
