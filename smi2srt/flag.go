package main

import (
	"flag"
	"log"
	"os"
)

var (
	flagOut = flag.String("o", "", "output srt file name")
	dest    *os.File
)

func init() {
	flag.Parse()

	if *flagOut != "" && flag.NArg() > 1 {
		log.Printf("smi2srt: flag -o will be ignored in multiple input")
	}

	if *flagOut != "" {
		if d, err := os.Create(*flagOut); err != nil {
			panic(err)
		} else {
			dest = d
		}
	} else {
		dest = os.Stdout
	}
}
