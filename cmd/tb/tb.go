package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/suapapa/tools/internal/tomboy"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log.SetPrefix("tb: ")
	log.SetFlags(0)

	dropboxPath := dropboxDir()

	tomboyPath := filepath.Join(dropboxPath, "Tomboy")
	if _, err := os.Stat(tomboyPath); os.IsNotExist(err) {
		return fmt.Errorf("tomboy path does not exist: %s", tomboyPath)
	}

	manifestPath := filepath.Join(tomboyPath, "manifest.xml")
	m, err := tomboy.ReadManifestFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	fmt.Printf("Last sync: %v\n", m.Date)
	fmt.Printf("Revision: %d\n", m.Revision)
	fmt.Printf("Notes: %d\n", len(m.Items))

	return nil
}
