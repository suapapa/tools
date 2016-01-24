// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
)

type table [][]string

func newTable(h []string) *table {
	t := new(table)
	t.AppendRows(h)
	return t
}

func (t *table) AppendRows(r ...[]string) {
	*t = append(*t, r...)
}

func (t *table) writeInCSV(w io.Writer) {
	for _, r := range *t {
		for i, c := range r {
			if i != 0 {
				w.Write([]byte{','})
			}
			w.Write([]byte(c))
		}
		w.Write([]byte{'\n'})
	}
}

func (t *table) writeInMd(w io.Writer) {
	colLens := make([]int, len((*t)[0]))
	for _, r := range *t {
		for i, c := range r {
			if colLens[i] < len(c) {
				colLens[i] = len(c)
			}
		}
	}

	fmtStrs := make([]string, len(colLens))
	for i, l := range colLens {
		fmtStrs[i] = " %" + fmt.Sprintf("% ds", l) + " "
	}

	// write header
	fmt.Fprint(w, "|")
	for i, c := range (*t)[0] {
		fmt.Fprintf(w, fmtStrs[i], c)
		fmt.Fprint(w, "|")
	}
	fmt.Fprintln(w)

	// write seperator
	fmt.Fprint(w, "|")
	for _, l := range colLens {
		for i := 0; i < l+2; i++ {
			fmt.Fprint(w, "-")
		}
		fmt.Fprint(w, "|")
	}
	fmt.Fprintln(w)

	// write value lines
	for _, r := range (*t)[1:] {
		fmt.Fprint(w, "|")
		for i, c := range r {
			fmt.Fprintf(w, fmtStrs[i], c)
			fmt.Fprint(w, "|")
		}
		fmt.Fprintln(w)
	}
}
