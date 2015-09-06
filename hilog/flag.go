// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	dumpLSP         bool
	listBundledLSPs bool
}

func setupFlags(opts *Options) *flag.FlagSet {
	prgName := os.Args[0]
	fs := flag.NewFlagSet(prgName, flag.ExitOnError)
	fs.BoolVar(&opts.dumpLSP, "p", false, "dump given LSP in JSON")
	fs.BoolVar(&opts.listBundledLSPs, "l", false, "list up bundled LSPs")

	fs.Usage = func() {
		fmt.Printf("Usage: %s [options] <LogSchemePack> command...\n",
			prgName)
		fs.PrintDefaults()
	}

	return fs
}

func verifyFlags(opts *Options, fs *flag.FlagSet) {
	if !opts.listBundledLSPs && len(fs.Args()) == 0 {
		fs.Usage()
		os.Exit(1)
	}
}

func parseFlags() (Options, []string) {
	var opts Options
	fs := setupFlags(&opts)
	fs.Parse(os.Args[1:])
	verifyFlags(&opts, fs)
	return opts, fs.Args()
}
