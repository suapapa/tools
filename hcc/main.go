package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

var flagGet bool
var flagSet string

func init() {
	flag.BoolVar(&flagGet, "g", false, "Get RTC time from HangulClock")
	flag.StringVar(&flagSet, "s", "now", "Set RTC time from HangulClock")
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("usage: hcc options serialport")
		os.Exit(-1)
	}

	port := flag.Arg(0)
	log.Println("Opening serial port", port)
	ser, err := serial.OpenPort(&serial.Config{
		Name:        port,
		Baud:        57600,
		ReadTimeout: 500 * time.Millisecond,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Please wait till flushing hello msgs.")
	time.Sleep(4 * time.Second)
	ser.Flush()

	readBuff := make([]byte, 100)

	switch {
	case flagGet == true:
		log.Println("Get time from HangulClock")
		fmt.Fprintf(ser, "#G\r\n")
		// 15:15:44 OK\r\n
		_, err := ser.Read(readBuff)
		if err != nil {
			log.Fatal("Fail to Get", err)
		}
		os.Stdout.WriteString("> ")
		os.Stdout.Write(readBuff)
	case flagSet == "now":
		log.Println("Set time to HangulClock")
		t := time.Now()
		log.Println("Set to system time,", t)
		h, m, s := t.Hour(), t.Minute(), t.Second()
		fmt.Fprintf(ser, "#S%02d%02d%02d\r\n", h, m, s)
		_, err := ser.Read(readBuff)
		if err != nil {
			log.Fatal("Fail to Get", err)
		}
		os.Stdout.WriteString("> ")
		os.Stdout.Write(readBuff)
	case flagSet != "now":
		log.Println("Set time to HangulClock")
		var h, m, s int
		fmt.Sscanf(flagSet, "%02d%02d%02d", &h, &m, &s)
		fmt.Fprintf(ser, "#S%02d%02d%02d\r\n", h, m, s)
		_, err := ser.Read(readBuff)
		if err != nil {
			log.Fatal("Fail to Get", err)
		}
		os.Stdout.WriteString("> ")
		os.Stdout.Write(readBuff)
	}
	log.Println("All done")
}
