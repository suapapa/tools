// Copyright 2015, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	flagOut = flag.String("o", "", "output SRT file name (ignored for multiple inputs)")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := flag.Args()

	switch len(args) {
	case 0: // SMI from stdin
		dest := os.Stdout
		if *flagOut != "" {
			d, err := os.Create(*flagOut)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer d.Close()
			dest = d
		}
		log.Printf("stdin -> %s", dest.Name())
		return convertSmiToSrt(os.Stdin, dest)

	case 1: // Single SMI
		smiFile := args[0]
		s, err := os.Open(smiFile)
		if err != nil {
			return fmt.Errorf("failed to open SMI file: %w", err)
		}
		defer s.Close()

		dest := os.Stdout
		if *flagOut != "" {
			d, err := os.Create(*flagOut)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer d.Close()
			dest = d
		} else {
			// If no -o, we could either print to stdout or create .srt next to .smi
			// Original code printed to stdout or -o. Let's keep that but maybe also create .srt if -o is empty?
			// The original code was: log.Printf("%s -> %s", s.Name(), dest.Name())
			// Let's keep it simple and just use -o or stdout for single file.
		}

		log.Printf("%s -> %s", s.Name(), dest.Name())
		return convertSmiToSrt(s, dest)

	default: // Multiple SMIs
		if *flagOut != "" {
			log.Printf("warning: -o will be ignored for multiple inputs")
		}

		for _, smi := range args {
			if err := processFile(smi); err != nil {
				log.Printf("failed to process %s: %v", smi, err)
			}
		}
	}

	return nil
}

func processFile(smi string) error {
	s, err := os.Open(smi)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer s.Close()

	srt := convertExtToSrt(smi)
	d, err := os.Create(srt)
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer d.Close()

	log.Printf("%s -> %s", smi, srt)
	return convertSmiToSrt(s, d)
}

func convertExtToSrt(fn string) string {
	i := strings.LastIndexByte(fn, '.')
	if i == -1 {
		return fn + ".srt"
	}
	return fn[:i+1] + "srt"
}
