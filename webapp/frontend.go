package webapp

import (
	_ "embed"
	"log"
	"net/http"
)

type indexPage []byte

func (ip indexPage) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(ip)
	if err != nil {
		log.Println(err)
	}
}

//go:embed index.gohtml
var indexPageContent indexPage
