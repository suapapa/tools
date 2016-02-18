// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

const defaultTimeFmt = "Mon, 2006-01-02 15:04:05 MST -0700"

func main() {
	err := loadDB()
	if err != nil {
		fmt.Println("failed to load DB", err)
		os.Exit(-1)
	}

	if len(os.Args) == 1 {
		printTimes(defaultTimeFmt)
		os.Exit(0)
	}

	c := cli.NewCLI("wdate", "1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"add": func() (cli.Command, error) {
			var cmd addCmd
			return cmd, nil
		},
		"remove": func() (cli.Command, error) {
			var cmd removeCmd
			return cmd, nil
		},
		"fmt": func() (cli.Command, error) {
			var cmd fmtCmd
			return cmd, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
