// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	opt  options
	args []string
)

func main() {
	opt, args = parseFlags()

	if len(args) < 2 {
		opt.usage()
		os.Exit(-1)
	}

	csvScheme := args[0]
	p, err := parseCsvFile(csvScheme)
	if err != nil {
		panic(err)
	}

	header := []string{"value"}
	header = append(header, p.Strings()...)
	table := newTable(header)

	for _, v := range args[1:] {
		row := []string{v}

		valueStr := strings.TrimPrefix(v, "0x")
		value, err := strconv.ParseUint(valueStr, 16, 0)
		if err != nil {
			panic(err)
		}

		var fmtStr string
		for _, f := range p {
			switch opt.base {
			case 2:
				fmtStr = "%" + fmt.Sprintf("0%db", f.cnt) + "b"
			case 10:
				fmtStr = "%d"
			case 16:
				fmtStr = "0x%" + fmt.Sprintf("0%dx", (f.cnt+3)/4)
			default:
				panic(fmt.Sprintf("unknown base %d", opt.base))
			}
			row = append(row, fmt.Sprintf(fmtStr, f.Value(value)))
		}

		table.AppendRows(row)
	}

	table.Render(os.Stdout)
}
