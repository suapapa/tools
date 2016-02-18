package main

import (
	"fmt"
	"time"
)

func printTimes(timeFmt string) {
	now := time.Now()
	fmt.Println(now.Format(timeFmt))
	for name, offset := range locs {
		loc := time.FixedZone(name, offset)
		fmt.Println(now.In(loc).Format(timeFmt))
	}
}
