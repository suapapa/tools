package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
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
	port := flag.Int("p", 8080, "listen port")
	count := flag.Int("c", 1, "number of downloads before exit")
	flag.Parse()

	serveFN := flag.Arg(0)
	if serveFN == "" {
		return fmt.Errorf("file to serve is required")
	}

	ip, err := netutil.FirstExternalIP()
	if err != nil {
		ip = "localhost"
	}

	log.Printf("Listen to http://%s:%d/%s", ip, *port, filepath.Base(serveFN))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	downloadsRemaining := *count
	var mu sync.Mutex

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", *port),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Access from %s", r.RemoteAddr)

		f, err := os.Open(serveFN)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(serveFN)))

		if _, err := io.Copy(w, f); err != nil {
			log.Printf("Copy failed: %v", err)
			return
		}

		mu.Lock()
		downloadsRemaining--
		if downloadsRemaining <= 0 {
			stop()
		}
		mu.Unlock()
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe failed: %v", err)
		}
		wg.Done()
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	wg.Wait()
	log.Println("Bye~")
	return nil
}
