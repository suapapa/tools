package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

var (
	flagOutput = flag.String("o", "out.png", "output PNG name")
	flagSize   = flag.Int("s", 512, "QR code size")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Usage: %s [options] <text>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := qrcode.WriteFile(flag.Arg(0), qrcode.Medium, *flagSize, *flagOutput); err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate QR code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("QR code saved to %s\n", *flagOutput)
}
