package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age: 10886400; includeSubdomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		fs.ServeHTTP(w, r)
	})

	addr := ":0"
	if value, ok := os.LookupEnv("ADDR"); ok {
		addr = value
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using port: %d", listener.Addr().(*net.TCPAddr).Port)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal(err)
	}
}
