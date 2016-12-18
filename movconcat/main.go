// Copyright 2016, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"
)

const movRePtn = `(\d\d\d\d_\d\d\d\d_\d\d\d\d\d\d)_\d\d\d.MOV`
const movTimeForm = "2006_0102_150405"

func main() {
	// list up MOVs
	// root := filepath.Join(os.Getenv("HOME"), "video/blackbox/raw")
	root := ""
	log.Println("searching MOVs from", root, "...")

	files, err := filepath.Glob(filepath.Join(root, "*.MOV"))
	panicIfErr(err)

	sort.Strings(files)
	log.Println(len(files), "MOVs found.")

	timeForm := regexp.MustCompile(movRePtn)

	// group MOVs by time
	var lastT time.Time
	var lastStartTime string
	movs := make(map[string][]string)
	for _, m := range files {
		// fmt.Println(filepath.Base(movs[0]))
		matchs := timeForm.FindStringSubmatch(m)
		if len(matchs) != 2 {
			log.Println("Skip", m)
		}

		currT, err2 := time.Parse(movTimeForm, matchs[1])
		panicIfErr(err2)

		if currT.Sub(lastT) > (10*time.Minute + 5*time.Second) {
			lastStartTime = matchs[1]
			log.Println("new recoding begins from", lastStartTime)
		}

		movs[lastStartTime] = append(movs[lastStartTime], m)

		lastT = currT
	}

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
					if !*flagIntermedeateFiles {
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
		if err != nil {
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
