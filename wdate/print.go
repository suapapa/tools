package main

import (
	"fmt"
	"time"
)

func printTimes(timeFmt string, t time.Time) {
	// t := time.Now()
	fmt.Printf("%s%d: %s\n",
		"w", weekCount(t.YearDay()), t.Format(timeFmt))
	for name, offset := range locs {
		locT := t.In(time.FixedZone(name, offset))
		fmt.Printf("%s%d: %s\n",
			"w", weekCount(locT.YearDay()), locT.Format(timeFmt))
	}
}

func weekCount(yd int) int {
	return yd/7 + 1
}
