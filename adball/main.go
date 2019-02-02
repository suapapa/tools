package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	devices := runAdbDevices()

	if len(os.Args) < 2 {
		fmt.Println("attached devices:")
		fmt.Println(devices)
		os.Exit(-1)
	}

	for _, serial := range devices {
		singleCmd := fmt.Sprintf("adb -s %s ", serial) + strings.Join(os.Args[1:], " ")
		out, err := runCommandSync(singleCmd)
		if err != nil {
			log.Fatalf("faile to run %s: %s", singleCmd, err)
		}
		fmt.Println(out)
	}
}
