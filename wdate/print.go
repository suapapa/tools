package main

import (
	"fmt"
	"time"
)

func printTimes(timeFmt string, t time.Time) {
	// t := time.Now()
	fmt.Printf("%s%d: %s\n",
		"w", weekCount(t), t.Format(timeFmt))
	for name, offset := range locs {
		locT := t.In(time.FixedZone(name, offset))
		fmt.Printf("%s%d: %s\n",
			"w", weekCount(locT), locT.Format(timeFmt))
	}
}

func weekCount(t time.Time) int {
	yd := t.YearDay()
	wd := t.Weekday()
	return yd/7 + 1
}
