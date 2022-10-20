package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/klauspost/compress/gzhttp"
)

var buildTime string

//go:embed static/*
var staticFiles embed.FS

func main() {
	// If no buildTime was set during compilation we use the
	// current time at startup.
	if buildTime == "" {
		buildTime = time.Now().UTC().Format(time.RFC1123Z)
	}

	log.Printf("Build time: %s", buildTime)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	webrootFS := http.FS(staticFS)

	fs := gzhttp.GzipHandler(http.FileServer(webrootFS))

	// Serve static files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubdomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Date", buildTime)

		// Don't list directories.
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

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
