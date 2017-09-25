package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	y := os.Args[1]
	currY := time.Now().Year()
	fromY, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}

	fmt.Println(currY - fromY)
}
