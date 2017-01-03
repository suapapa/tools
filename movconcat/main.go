// Copyright 2016, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func main() {
	// list up MOVs
	// root := filepath.Join(os.Getenv("HOME"), "video/blackbox/raw")
	root := ""
	log.Println("searching MOVs from", root, "...")

	files, err := filepath.Glob(filepath.Join(root, "*.MOV"))
	panicIfErr(err)

	sort.Strings(files)
	log.Println(len(files), "MOVs found.")

	movs := sjChapter(files)

	// concat MOVs
	if *flagJobs <= 0 {
		*flagJobs = 1
	}

	var wg sync.WaitGroup
	wg.Add(*flagJobs)

	type Clips struct {
		k string
		v []string
	}

	workC := make(chan *Clips)
	errC := make(chan error)
	ctx, cancle := context.WithCancel(context.Background())

	for i := 0; i < *flagJobs; i++ {
		// create Workers
		go func(id int, ctx context.Context) {
			log.Printf("worker %d start\n", id)
			defer wg.Done()
		loop:
			for {
				select {
				case <-ctx.Done():
					log.Printf("worker %d finish\n", id)
					break loop
				case c := <-workC:
					runFFmpeg(c.k, c.v)
					if !*flagDryrun && !*flagIntermedeateFiles {
						for _, f := range c.v {
							os.Remove(f)
						}
					}
				}
			}
			errC <- ctx.Err()
		}(i, ctx)
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
