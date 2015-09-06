package main

import (
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("usage: " + os.Args[0] + " dump_name w h")
		os.Exit(-1)
	}

	fn := os.Args[1]
	w, err := strconv.Atoi(os.Args[2])
	exitIf(err)
	h, err := strconv.Atoi(os.Args[3])
	exitIf(err)

	dump, err := os.Open(fn)
	exitIf(err)
	defer dump.Close()

	gray := image.NewGray(image.Rect(0, 0, w, h))
	n, err := dump.Read(gray.Pix)
	exitIf(err)

	if n != w*h {
		exitIf(errors.New("dump is too small"))
	}

	toimg, err := os.Create(fn + ".png")
	exitIf(err)
	defer toimg.Close()

	err = png.Encode(toimg, gray)
	exitIf(err)
}

func exitIf(err error) {
	if err != nil {
		_, _, ln, _ := runtime.Caller(1)
		log.Println("err:", ln, err)
		os.Exit(-1)
	}
}
