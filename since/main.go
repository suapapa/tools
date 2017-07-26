package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	units "github.com/docker/go-units"
)

const (
	timeForm = "2006-01-02"
)

func main() {
	toTime := flag.String("t", "now", "to time")
	flag.Parse()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(-1)
	}

	tFrom, err := time.Parse(timeForm, flag.Arg(0))
	if err != nil {
		log.Fatal("failed to parse time string:", err)
	}

	var tTo time.Time
	if *toTime == "now" {
		tTo = time.Now()
	} else {
		tTo, err = time.Parse(timeForm, *toTime)
		if err != nil {
			log.Fatal("failed to parse time string:", err)
		}
	}

	fmt.Println(units.HumanDuration(tTo.Sub(tFrom)))
}

func printUsage() {
	fmt.Println("usage:", os.Args[0], "2017-05-10")
}
