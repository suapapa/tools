// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Scale represents,
type Scale struct {
	n    int
	name []string
}

var (
	flagScale string
)

func init() {
	flag.StringVar(&flagScale, "s", "ko", "Select scale in ko, en, bin or all")
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		fmt.Printf("usage: %s Numbers\n", os.Args[0])
		os.Exit(-1)
	}

	scales := []*Scale{}

	switch flagScale {
	case "ko":
		scales = append(scales, &koScale)
	case "en":
		scales = append(scales, &enScale)
	case "bin":
		scales = append(scales, &binScale)
	case "all":
		scales = append(scales, &koScale, &enScale, &binScale)
	default:
		fmt.Fprint(os.Stderr, "unsupported scale. use ko, en, bin or all\n")
		os.Exit(-1)
	}

	for _, n := range flag.Args() {
		for _, scale := range scales {
			ns, err := convertForHuman(n, scale)
			if err != nil {
				// TODO: beautiful usage
				fmt.Fprintln(os.Stderr, err)
				flag.PrintDefaults()
				os.Exit(-1)
			}
			fmt.Print(ns, " ")
		}
		fmt.Println()
	}

}

func isNum(s string) error {
	for _, c := range s {
		if !(('0' <= c) && (c <= '9')) {
			return errors.New("input is not number")
		}
	}
	return nil
}

func convertForHuman(n string, s *Scale) (string, error) {
	err := isNum(n)
	if err != nil {
		return "", err
	}

	l := len(n)

	mIdx := l / s.n
	if mIdx >= len(s.name) {
		return "", errors.New("too big to current scale")
	}

	if l%s.n == 0 {
		mIdx--
	}
	metric := s.name[mIdx]

	hTo := l - (mIdx * s.n)
	head := n[:hTo]

	return head + metric, nil
}
