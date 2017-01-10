// Copyright 2015, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	switch flag.NArg() {
	case 0: // smi from stdin
		log.Printf("%s -> %s", os.Stdin.Name(), dest.Name())
		if err := smi2srt(os.Stdin, dest); err != nil {
			panic(err)
		}

	case 1: // single smi
		s, err := os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}

		log.Printf("%s -> %s", s.Name(), dest.Name())
		if err := smi2srt(s, dest); err != nil {
			panic(err)
		}

		s.Close()

	default: // multiple smi
		for _, smi := range flag.Args() {
			s, err := os.Open(smi)
			if err != nil {
				panic(err)
			}

			d, err := os.Create(convertExtToSrt(smi))
			if err != nil {
				panic(err)
			}

			log.Printf("%s -> %s", s.Name(), d.Name())
			if err := smi2srt(s, d); err != nil {
				panic(err)
			}

			d.Close()
			s.Close()
		}
	}

	// exit
	dest.Close()
}

func convertExtToSrt(fn string) string {
	i := strings.LastIndexByte(fn, '.')

	return fn[:i+1] + "srt"
}
