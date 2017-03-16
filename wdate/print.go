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
	// TODO: just int(t.Weekday()) OK?
	var offset int
	offset = int(t.Weekday())
	// offset := int(fDate.Weekday()) - 1

	yd := t.YearDay()

	r := ((yd + offset) / 7)
	// log.Println("weekCount:", t, yd, offset, r)
	// log.Println()
	return r
}
