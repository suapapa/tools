// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"log"
)

var errInvalidOffset = errors.New("Invalid time offset format")

type addCmd struct{}

func (f addCmd) Synopsis() string {
	return "add a timezone"
}

func (f addCmd) Help() string {
	return "add a timezone in form of 'TZNAME +0700'"
}

func (f addCmd) Run(args []string) int {
	if len(args) != 2 {
		log.Println("fail to add a timezone")
		return -1
	}
	name := args[0]
	offsetStr := args[1]
	offset, err := timeOffsetStrToInt(offsetStr)
	if err != nil {
		log.Println("fail to convert offset:", err)
		return -1
	}

	err = addLocation(name, offset)
	if err != nil {
		log.Println("fail to add location:", err)
		return -1
	}

	return 0
}
