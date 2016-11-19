package main

import (
	"flag"
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

var (
	flagUseDocker               = flag.Bool("d", false, "use docker")
	flagParalleRun              = flag.Bool("p", false, "run cmd parallely")
	flagDeleteIntermedeateFiles = flag.Bool("d", false,
		"delete intermedeate files after finish concat")
)

func main() {
	flag.Parse()
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
	var wg sync.WaitGroup
	wg.Add(len(movs))
	for k, v := range movs {
		if *flagParalleRun {
			go func(k string, v []string) {
				runFFmpeg(k, v)
				if *flagDeleteIntermedeateFiles {
					for _, f := range v {
						os.Remove(f)
					}
				}
				wg.Done()
			}(k, v)
		} else {
			runFFmpeg(k, v)
			if *flagDeleteIntermedeateFiles {
				for _, f := range v {
					os.Remove(f)
				}
			}
			wg.Done()
		}
	}

	wg.Wait()
	log.Println("all done!")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
