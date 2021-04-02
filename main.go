package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	webrootFS := http.FS(staticFS)

	fs := http.FileServer(webrootFS)

	// Serve static files
	http.Handle("/", fs)

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
