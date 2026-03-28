package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	devices, err := runAdbDevices(ctx)
	if err != nil {
		return err
	}

	if len(os.Args) < 2 {
		fmt.Println("attached devices:")
		for _, d := range devices {
			fmt.Printf(" - %s\n", d)
		}
		return nil
	}

	for _, serial := range devices {
		args := append([]string{"-s", serial}, os.Args[1:]...)
		out, err := runCommand(ctx, "adb", args...)
		if err != nil {
			log.Printf("failed to run adb %v for serial %s: %v", os.Args[1:], serial, err)
			continue
		}
		fmt.Println(out)
	}

	return nil
}
