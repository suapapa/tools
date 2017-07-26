package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
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

	t, err := time.Parse(timeForm, flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	var d time.Duration
	if *toTime == "now" {
		d = time.Since(t)
	} else {
		tt, err := time.Parse(timeForm, *toTime)
		if err != nil {
			log.Fatal(err)
		}
		d = tt.Sub(t)
	}

	// TODO: it prints duration by hours. Need human readable form
	fmt.Println(d)
}

func printUsage() {
	fmt.Println("usage:", os.Args[0], "2017-05-10")
}
