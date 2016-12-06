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

	wCh := make(chan *Clips)
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	for i := 0; i < *flagJobs; i++ {
		// create Workers
		go func(ctx context.Context, ch chan *Clips) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case c := <-ch:
					runFFmpeg(c.k, c.v)
					if *flagDeleteIntermedeateFiles {
						for _, f := range c.v {
							os.Remove(f)
						}
					}
				}
			}
		}(ctx, wCh)
	}

	for k, v := range movs {
		wCh <- &Clips{k: k, v: v}
	}

	wg.Wait()
	log.Println("all done!")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
