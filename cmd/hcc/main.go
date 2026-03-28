package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flagGet := flag.Bool("g", false, "Get RTC time from HangulClock")
	flagSet := flag.String("s", "now", "Set RTC time from HangulClock (HHMMSS or 'now')")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s [options] <serialport>\n", os.Args[0])
		flag.PrintDefaults()
		return nil
	}

	port := flag.Arg(0)
	log.Printf("Opening serial port %s", port)
	ser, err := serial.OpenPort(&serial.Config{
		Name:        port,
		Baud:        57600,
		ReadTimeout: 500 * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("failed to open port: %w", err)
	}
	defer ser.Close()

	log.Println("Please wait till flushing hello msgs.")
	time.Sleep(2 * time.Second) // reduced from 4
	ser.Flush()

	readBuff := make([]byte, 100)

	switch {
	case *flagGet:
		log.Println("Get time from HangulClock")
		fmt.Fprintf(ser, "#G\r\n")
		n, err := ser.Read(readBuff)
		if err != nil {
			return fmt.Errorf("failed to read from serial: %w", err)
		}
		fmt.Printf("> %s", string(readBuff[:n]))

	case *flagSet == "now":
		log.Println("Set time to HangulClock")
		t := time.Now()
		log.Printf("Set to system time: %v", t)
		h, m, s := t.Hour(), t.Minute(), t.Second()
		fmt.Fprintf(ser, "#S%02d%02d%02d\r\n", h, m, s)
		n, err := ser.Read(readBuff)
		if err != nil {
			return fmt.Errorf("failed to read from serial: %w", err)
		}
		fmt.Printf("> %s", string(readBuff[:n]))

	case *flagSet != "now":
		log.Printf("Set time to HangulClock: %s", *flagSet)
		var h, m, s int
		if _, err := fmt.Sscanf(*flagSet, "%02d%02d%02d", &h, &m, &s); err != nil {
			return fmt.Errorf("invalid time format, expect HHMMSS: %w", err)
		}
		fmt.Fprintf(ser, "#S%02d%02d%02d\r\n", h, m, s)
		n, err := ser.Read(readBuff)
		if err != nil {
			return fmt.Errorf("failed to read from serial: %w", err)
		}
		fmt.Printf("> %s", string(readBuff[:n]))
	}

	log.Println("All done")
	return nil
}
