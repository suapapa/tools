package main

import (
	"flag"
	"runtime"
)

var (
	flagUseDocker               = flag.Bool("d", false, "use docker")
	flagJobs                    = flag.Int("j", runtime.NumCPU(), "parallel jobs")
	flagDeleteIntermedeateFiles = flag.Bool("i", false,
		"delete intermedeate files after finish concat")
)

func init() {
	flag.Parse()
}
