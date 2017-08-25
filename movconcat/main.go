// Copyright 2016, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime/trace"
)

func main() {
	// trace
	if *flagTrace {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	// list up MOVs
	// root := filepath.Join(os.Getenv("HOME"), "video/blackbox/raw")
	root := ""
	log.Println("searching MOVs from", root, "...")

	var movs map[string][]string

	// searching SJCam videos
	files, err := filepath.Glob(filepath.Join(root, "*.MOV"))
	panicIfErr(err)

	if len(files) != 0 {
		movs = sjChapter(files)

	} else {
		files, err = filepath.Glob(filepath.Join(root, "*.MP4"))
		panicIfErr(err)
		if len(files) != 0 {
			movs = gpChapter(files)
		}
	}

	if *flagJobs > len(movs) {
		*flagJobs = len(movs)
	}

	type Clips struct {
		k string
		v []string
	}

	workC := make(chan *Clips)
	errC := make(chan error)
	ctx, cancle := context.WithCancel(context.Background())

	for i := 0; i < *flagJobs; i++ {
		// create Workers
		go func(ctx context.Context, id int) {
			log.Printf("worker %d start\n", id)
		loop:
			for {
				select {
				case <-ctx.Done():
					log.Printf("worker %d finish\n", id)
					break loop
				case c := <-workC:
					err := runFFmpeg(c.k, c.v)
					if err != nil {
						panic(err)
					}
					if !*flagDryrun && !*flagIntermedeateFiles {
						for _, f := range c.v {
							os.Remove(f)
						}
					}
				}
			}
			errC <- ctx.Err()
		}(ctx, i)
	}

	for k, v := range movs {
		workC <- &Clips{k: k, v: v}
	}

	// finish workers
	cancle()

	for i := 0; i < *flagJobs; i++ {
		err := <-errC
		if err != context.Canceled && err != nil {
			panic(err)
		}
	}

	log.Println("all done!")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
