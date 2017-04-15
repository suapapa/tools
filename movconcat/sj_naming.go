package main

import (
	"log"
	"regexp"
	"sort"
	"time"
)

const movRePtn = `(\d\d\d\d_\d\d\d\d_\d\d\d\d\d\d)_\d\d\d.MOV`
const movEveRePtn = `(\d\d\d\d_\d\d\d\d_\d\d\d\d\d\d)_\d\d\d_EVE.MOV`
const movTimeForm = "2006_0102_150405"

var (
	timeForm    = regexp.MustCompile(movRePtn)
	timeEveForm = regexp.MustCompile(movEveRePtn)
)

// sjChapter makes filename map of starttime
func sjChapter(files []string) map[string][]string {
	sort.Strings(files)
	log.Println(len(files), "MOVs found.")

	// group MOVs by time
	var lastT time.Time
	var lastStartTime string
	movs := make(map[string][]string)
loop:
	for _, m := range files {
		// fmt.Println(filepath.Base(movs[0]))
		matchs := timeForm.FindStringSubmatch(m)
		if len(matchs) != 2 {
			matchs := timeEveForm.FindStringSubmatch(m)
			if len(matchs) != 2 {
				log.Println("Skip", m)
				continue loop
			}
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
	return movs
}
