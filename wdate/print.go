package main

import (
	"fmt"
	"time"
)

func printTimes(timeFmt string) {
	now := time.Now()
	fmt.Printf("%s%d: %s\n", "w", now.YearDay()/7, now.Format(timeFmt))
	for name, offset := range locs {
		locNow := now.In(time.FixedZone(name, offset))
		fmt.Printf("%s%d: %s\n", "w", locNow.YearDay()/7, locNow.Format(timeFmt))
	}
}
