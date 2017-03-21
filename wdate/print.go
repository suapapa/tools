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
	ft := firstdayOfTheYear(t)
	offset := 7 - weekdayOffsetFromMonday(ft.Weekday())

	yd := t.YearDay() - offset
	r := (yd-1)/7 + 1

	// log.Printf("offset:%d yd:%d r:%d\n",
	// 	offset, yd, r)

	return r
}

func weekdayOffsetFromMonday(d time.Weekday) int {
	if d == time.Sunday {
		return 6
	}

	return int(d) - 1
}

func firstdayOfTheYear(t time.Time) time.Time {
	f, err := time.Parse("2006-01-02", fmt.Sprintf("%04d-01-01", t.Year()))
	if err != nil {
		panic(err)
	}

	return f
}
