// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runFFmpeg(k string, v []string) {
	o := k + ".mov"

	if len(v) == 1 {
		os.Rename(v[0], o)
		return
	}

	tmp, err := os.Create(k + ".list")
	panicIfErr(err)
	defer os.Remove(tmp.Name())

	for _, f := range v {
		fmt.Fprintln(tmp, "file", f)
	}
	tmp.Close()

	var cmd *exec.Cmd
	if *flagUseDocker {
		// cf) https://github.com/jrottenberg/ffmpeg
		// $ sudo docker run -v $PWD:/opt/data --rm \
		// jrottenberg/ffmpeg -f concat \
		// -i /opt/data/files.list -c copy /opt/data/test.mov
		cmd = exec.Command("docker", "run", "--rm",
			"-v", os.Getenv("PWD")+":/opt/data",
			"jrottenberg/ffmpeg",
			"-f", "concat",
			"-i", filepath.Join("/opt/data", tmp.Name()),
			"-c", "copy",
			filepath.Join("/opt/data", o),
		)
	} else {
		cmd = exec.Command("ffmpeg",
			"-f", "concat",
			"-i", tmp.Name(),
			"-c", "copy",
			o,
		)
	}

	if *flagDeleteIntermedeateFiles == false {
		stdErr, err2 := os.Create(k + ".log")
		panicIfErr(err2)
		defer stdErr.Close()
		cmd.Stderr = stdErr
	}

	log.Println("concatting", o, "...")

	err = cmd.Start()
	panicIfErr(err)

	err = cmd.Wait()
	if err != nil {
		log.Printf("failed to concat %s with error: %v\n", o, err)
	} else {
		log.Println(o, "finished.")
	}
}
