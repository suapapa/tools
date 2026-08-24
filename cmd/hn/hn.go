// Copyright 2015, Homin Lee <i@homin.dev>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/suapapa/tools/pkg/scale"
)

var (
	flagScale = flag.String("s", "ko", "Select scale in ko, en, bin or all")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage: %s Numbers\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	scales := []scale.Scale{}

	switch *flagScale {
	case "ko":
		scales = append(scales, scale.Ko)
	case "en":
		scales = append(scales, scale.En)
	case "bin":
		scales = append(scales, scale.Bin)
	case "all":
		scales = append(scales, scale.Ko, scale.En, scale.Bin)
	default:
		fmt.Fprintf(os.Stderr, "unsupported scale %q. use ko, en, bin or all\n", *flagScale)
		os.Exit(1)
	}

	for _, n := range flag.Args() {
		for _, s := range scales {
			ns, err := scale.Format(n, s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				continue
			}
			fmt.Print(ns, " ")
		}
		fmt.Println()
	}
}
