package main

import (
	"fmt"
	"os"
	"time"

	"github.com/suapapa/tools/pref/iomon"
)

func main() {
	// TODO: bufsize in flag
	bufsize := 4096
	buf := make([]byte, bufsize)

	r := iomon.NewReadMon(os.Stdin)
	w := iomon.NewWriteMon(os.Stdout)

	for {
		rn, _ := r.Read(buf)
		wn, _ := w.Write(buf)

		if rn == wn && rn != bufsize {
			break
		}

		rnc, rtc := r.Check()
		wnc, wtc := w.Check()

		// TODO: worry overflow
		rbps := rnc * (uint64(time.Second) / 1024) / uint64(time.Since(rtc))
		wbps := wnc * (uint64(time.Second) / 1024) / uint64(time.Since(wtc))

		fmt.Fprintf(os.Stderr, "Read: %d kbps ", rbps)
		fmt.Fprintf(os.Stderr, "Write: %d kbps\n", wbps)
	}
}
