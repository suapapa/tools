package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	units "github.com/docker/go-units"
)

const timeForm = "2006-01-02"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	toTime := flag.String("t", "now", "to time (YYYY-MM-DD or 'now')")
	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		return nil
	}

	tFrom, err := time.Parse(timeForm, flag.Arg(0))
	if err != nil {
		return fmt.Errorf("failed to parse from time: %w", err)
	}

	var tTo time.Time
	if *toTime == "now" {
		tTo = time.Now()
	} else {
		tTo, err = time.Parse(timeForm, *toTime)
		if err != nil {
			return fmt.Errorf("failed to parse to time: %w", err)
		}
	}

	fmt.Println(units.HumanDuration(tTo.Sub(tFrom)))
	return nil
}

func printUsage() {
	fmt.Printf("Usage: %s [options] <YYYY-MM-DD>\n", os.Args[0])
	flag.PrintDefaults()
}
