package webapp

import (
	"embed"
	_ "embed"
	"io/fs"
	"net/http"
)

//go:embed site/*
var site embed.FS

var ui http.Handler

func init() {
	s, err := fs.Sub(site, "site")
	if err != nil {
		panic("can't access the ui files")
	}

	ui = http.FileServer(http.FS(s))
}
