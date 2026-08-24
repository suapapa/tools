package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var (
	flagBox        = flag.Uint("b", 800, "max width/height")
	flagResampling = flag.String("r", "NearestNeighbor",
		"resampler: NearestNeighbor, Bilinear, Bicubic, MitchellNetravali, Lanczos2, Lanczos3")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	imgs := flag.Args()

	if len(imgs) == 0 {
		fmt.Printf("Usage: %s [options] imagefile ...\n", os.Args[0])
		flag.PrintDefaults()
		return nil
	}

	for _, i := range imgs {
		if err := processImage(i); err != nil {
			log.Printf("failed to process %s: %v", i, err)
		}
	}

	return nil
}

func processImage(fn string) error {
	img, err := openImage(fn)
	if err != nil {
		return err
	}

	t := resize.Thumbnail(*flagBox, *flagBox, img, resampler(*flagResampling))

	outFn := outName(fn)
	out, err := os.Create(outFn)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, t, &jpeg.Options{Quality: 85}); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	fmt.Printf("%s -> %s\n", fn, outFn)
	return nil
}

func openImage(fn string) (image.Image, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	const maxPixels = 64 << 20
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxPixels {
		return nil, fmt.Errorf("image too large: %dx%d", cfg.Width, cfg.Height)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to rewind image: %w", err)
	}

	i, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return i, nil
}

func resampler(n string) resize.InterpolationFunction {
	switch n {
	case "Bilinear":
		return resize.Bilinear
	case "Bicubic":
		return resize.Bicubic
	case "MitchellNetravali":
		return resize.MitchellNetravali
	case "Lanczos2":
		return resize.Lanczos2
	case "Lanczos3":
		return resize.Lanczos3
	default:
		return resize.NearestNeighbor
	}
}

func outName(n string) string {
	ext := filepath.Ext(n)
	base := strings.TrimSuffix(n, ext)
	return fmt.Sprintf("%s_%d.jpg", base, *flagBox)
}
