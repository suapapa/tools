package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	dryRun := flag.Bool("n", false, "dry-run")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s [-n] zipfile ...\n", os.Args[0])
		return nil
	}

	euckrDec := korean.EUCKR.NewDecoder()

	for _, zipFilename := range flag.Args() {
		fmt.Printf("Unzipping %s ...\n", zipFilename)
		if err := convUnzip(zipFilename, euckrDec, *dryRun); err != nil {
			log.Printf("failed to unzip %s: %v", zipFilename, err)
		}
	}
	return nil
}

func convUnzip(fileName string, dec transform.Transformer, dryRun bool) error {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fixed, _, err := transform.String(dec, f.Name)
		if err != nil {
			log.Printf("failed to transform filename %s: %v", f.Name, err)
			fixed = f.Name // fallback
		}

		fmt.Println(fixed)

		if dryRun {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fixed, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fixed), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		if err := extractFile(f, fixed); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open zip file member: %w", err)
	}
	defer rc.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, rc); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}
	return nil
}
