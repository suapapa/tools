package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s birth_year\n", os.Args[0])
		os.Exit(1)
	}

	y := os.Args[1]
	currY := time.Now().Year()
	fromY, err := strconv.Atoi(y)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid year: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(currY - fromY)
}
