// https://suapapa.github.io/resume/

package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
)

var (
	flagOutput string
)

func main() {
	flag.StringVar(&flagOutput, "o", "out.png", "output png name")
	flag.Parse()

	png, err := qrcode.Encode(os.Args[1], qrcode.Medium, 512)
	if err != nil {
		log.Fatal(err)
	}

	pngReader := bytes.NewReader(png)

	w, err := os.Create(flagOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	io.Copy(w, pngReader)
}
