package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <dump_name> <width> <height>\n", os.Args[0])
		return nil
	}

	fn := os.Args[1]
	w, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return fmt.Errorf("invalid width: %w", err)
	}
	h, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return fmt.Errorf("invalid height: %w", err)
	}

	dump, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("failed to open dump file: %w", err)
	}
	defer dump.Close()

	gray := image.NewGray(image.Rect(0, 0, w, h))
	n, err := io.ReadFull(dump, gray.Pix)
	if err != nil {
		return fmt.Errorf("failed to read dump: %w", err)
	}

	if n != w*h {
		return fmt.Errorf("dump is too small: read %d bytes, expected %d", n, w*h)
	}

	outFn := fn + ".png"
	toimg, err := os.Create(outFn)
	if err != nil {
		return fmt.Errorf("failed to create PNG file: %w", err)
	}
	defer toimg.Close()

	if err := png.Encode(toimg, gray); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	fmt.Printf("Converted %s to %s (%dx%d)\n", fn, outFn, w, h)
	return nil
}
