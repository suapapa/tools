package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

var (
	flagBox        = flag.Uint("b", 800, "max width/height")
	flagResampling = flag.String("r", "NearestNeighbor",
		"resampler can be following:"+
			"\n\t"+
			"NearestNeighbor "+
			"Bilinear "+
			"Bicubic "+
			"MitchellNetravali "+
			"Lanczos2 "+
			"Lanczos3 "+
			"\n\t",
	)
)

func main() {
	flag.Parse()
	imgs := flag.Args()

	if len(imgs) == 0 {
		fmt.Printf("usage: %s imagefile ...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, i := range imgs {
		img, err := openImage(i)
		if err != nil {
			panic(err)
		}

		t := resize.Thumbnail(*flagBox, *flagBox, img, resampler(*flagResampling))
		out, err := os.Create(outName(i))
		if err != nil {
			panic(err)
		}

		err = png.Encode(out, t)
		if err != nil {
			panic(err)
		}
		out.Close()
	}
}

func openImage(fn string) (image.Image, error) {
	// log.Println(fn)
	// fn = strings.Replace(fn, "~", os.Getenv("HOME"), 1)

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func resampler(n string) resize.InterpolationFunction {
	switch n {
	// case "NearestNeighbor":
	// 	return resize.NearestNeighbor
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
	}
	return resize.NearestNeighbor
}

func outName(n string) string {
	ss := strings.Split(n, ".")
	return strings.Join(ss[:len(ss)-1], ".") +
		"_" + fmt.Sprint(*flagBox) +
		".png"
}
