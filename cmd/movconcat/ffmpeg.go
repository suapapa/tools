// Copyright 2016, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runFFmpeg(k string, v []string) error {
	if len(v) == 0 {
		return fmt.Errorf("noting to concat")
	}

	o := k + "." + ext(v[0])
	log.Printf("concatting %s(%d clips)\n", o, len(v))

	if *flagDryrun {
		log.Println(o, "finished. (dry run)")
		return nil
	}

	if len(v) == 1 {
		os.Rename(v[0], o)
		return nil
	}

	tmp, err := os.Create(k + ".list")
	if err != nil {
		return fmt.Errorf("ffmpeg: failed to create list: %v", err)
	}
	defer os.Remove(tmp.Name())

	for _, f := range v {
		line, err := concatListEntry(f)
		if err != nil {
			tmp.Close()
			return err
		}
		if _, err := fmt.Fprintln(tmp, line); err != nil {
			tmp.Close()
			return fmt.Errorf("ffmpeg: failed to write list: %v", err)
		}
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("ffmpeg: failed to sync list: %v", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("ffmpeg: failed to close list: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("ffmpeg: fail to get workdir: %v", err)
	}

	var cmd *exec.Cmd
	if *flagUseDocker {
		// cf) https://github.com/jrottenberg/ffmpeg
		// $ sudo docker run -v $PWD:/opt/data --rm \
		// jrottenberg/ffmpeg -f concat \
		// -i /opt/data/files.list -c copy /opt/data/test.mov
		cmd = exec.Command("docker", "run", "--rm",
			// "--user", // for not owner root for output files
			"-v", wd+":/opt/data",
			"-w", "/opt/data",
			"jrottenberg/ffmpeg",
			"-f", "concat",
			"-safe", "0",
			"-i", "/opt/data/"+tmp.Name(),
			"-c", "copy",
			"/opt/data/"+o,
		)
	} else {
		cmd = exec.Command("ffmpeg",
			"-f", "concat",
			"-safe", "0",
			"-i", tmp.Name(),
			"-c", "copy",
			o,
		)
	}

	if *flagIntermedeateFiles {
		stdErr, err2 := os.Create(k + ".log")
		if err2 != nil {
			return fmt.Errorf("ffmpeg: cannot create log: %v", err2)
		}
		defer stdErr.Close()
		cmd.Stderr = stdErr
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("ffmpeg: fail to start cmd: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("ffmpeg: failed to concat %v", err)
	}

	log.Println(o, "finished.")
	return nil
}

func concatListEntry(path string) (string, error) {
	if strings.ContainsAny(path, "\n\r") {
		return "", fmt.Errorf("ffmpeg: filename contains newline: %q", path)
	}
	escaped := strings.ReplaceAll(path, `'`, `'\''`)
	return "file '" + escaped + "'", nil
}

func ext(fn string) string {
	s := strings.Split(fn, ".")
	return strings.ToLower(s[1])
}
