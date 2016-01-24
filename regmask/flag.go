// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

type options struct {
	fmt         string
	base        int
	headerGuide bool
	outFile     string

	usage func()
}

func setupFlags(opts *options) *flag.FlagSet {
	prgName := os.Args[0]
	fs := flag.NewFlagSet(prgName, flag.ExitOnError)
	fs.StringVar(&opts.fmt, "f", "md", "output format")
	fs.IntVar(&opts.base, "b", 2, "output base")
	fs.BoolVar(&opts.headerGuide, "g", false, "print guide in header")
	fs.StringVar(&opts.outFile, "o", "", "write output to a file")

	fs.Usage = func() {
		fmt.Printf("Usage: %s [options] <scheme> value...\n",
			prgName)
		fs.PrintDefaults()
	}

	return fs
}

func verifyFlags(opts *options, fs *flag.FlagSet) {
	switch opts.base {
	case 2, 10, 16:
	default:
		fs.Usage()
		os.Exit(1)
	}

	switch opts.fmt {
	case "md", "csv":
	default:
		fs.Usage()
		os.Exit(1)
	}
}

func parseFlags() (options, []string) {
	var opts options
	fs := setupFlags(&opts)
	fs.Parse(os.Args[1:])
	opts.usage = fs.Usage
	verifyFlags(&opts, fs)
	return opts, fs.Args()
}
