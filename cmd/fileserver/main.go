package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/suapapa/tools/pkg/netutil"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fPort := flag.String("p", "8080", "server port")
	fDir := flag.String("d", ".", "directory to serve")
	flag.Parse()

	dir, err := filepath.Abs(*fDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	log.Printf("Serving %s ...", dir)

	ip, err := netutil.FirstExternalIP()
	if err != nil {
		log.Printf("warning: %v", err)
		ip = "localhost"
	}

	addr := fmt.Sprintf(":%s", *fPort)
	log.Printf("Listening on http://%s%s", ip, addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           http.FileServer(http.Dir(dir)),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
	return srv.ListenAndServe()
}
