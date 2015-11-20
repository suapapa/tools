package main

import (
	"flag"
	"os"

	"github.com/suapapa/go_subtitle"
)

var (
	flagOut string
)

func init() {
	flag.StringVar(&flagOut, "o", "", "srt file name")
}

func main() {
	flag.Parse()

	var w *os.File
	var err error
	if flagOut != "" {
		w, err = os.Create(flagOut)
		if err != nil {
			panic(err)
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	smiFileName := flag.Arg(0)
	s, err := os.Open(smiFileName)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// read smi
	b, err := subtitle.ReadSmi(s)
	if err != nil {
		panic(err)
	}

	err = subtitle.ExportToSrtFile(b, w)
	if err != nil {
		panic(err)
	}
}
