// Copyright 2016, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"runtime"
)

var (
	flagUseDocker         = flag.Bool("d", false, "use docker")
	flagJobs              = flag.Int("j", runtime.NumCPU(), "parallel jobs")
	flagIntermedeateFiles = flag.Bool("i", false,
		"don't delete intermedeate files after finish concat")
	flagDryrun = flag.Bool("n", false, "dry run")
)

func init() {
	flag.Parse()
}
