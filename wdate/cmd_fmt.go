// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
)

type fmtCmd struct{}

func (f fmtCmd) Synopsis() string {
	return "set time format"
}

func (f fmtCmd) Help() string {
	return fmt.Sprintf("show times in given format. (default = \"%s\")",
		defaultTimeFmt)
}

func (f fmtCmd) Run(args []string) int {
	if len(args) != 1 {
		log.Println("fail to fmt. give a time-format")
		return -1
	}

	fmt := args[0]
	printTimes(fmt)

	return 0
}
