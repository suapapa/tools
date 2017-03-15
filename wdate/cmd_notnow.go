// Copyright 2017, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"
)

// var errInvalidDate = errors.New("Invalid date format")

type notNowCmd struct{}

func (f notNowCmd) Synopsis() string {
	return "specify a date"
}

func (f notNowCmd) Help() string {
	return "specify a date in form of '2017-03-15'"
}

func (f notNowCmd) Run(args []string) int {
	if len(args) != 1 {
		log.Println("fail to specify a date")
		return -1
	}

	t, err := parseTimeStr(args[0])
	if err != nil {
		log.Println("wrong date format:", err)
		return -1
	}

	printTimes(defaultTimeFmt, t)

	return 0
}

func parseTimeStr(tStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", tStr)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
