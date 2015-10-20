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
	flag.StringVar(&flagScale, "s", "ko", "Select Scale")
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s Number\n", os.Args[0])
		os.Exit(-1)
	}

	scale := &koScale
	switch flagScale {
	case "ko":
		scale = &koScale
	case "en":
		scale = &enScale
	case "bin":
		scale = &binScale
	default:
		fmt.Fprint(os.Stderr, "unsupported scale. use ko, en or bin\n")
		os.Exit(-1)
	}

	n := flag.Arg(0)
	ns, err := convertForHuman(n, scale)
	if err != nil {
		// TODO: beautiful usage
		fmt.Fprintln(os.Stderr, err)
		flag.PrintDefaults()
		os.Exit(-1)
	}

	fmt.Println(ns)
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
