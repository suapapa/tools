package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	j, _ := getID(116)
	// fmt.Println(j)

	var w Success
	err := json.Unmarshal(j, &w)
	if err != nil {
		var f Fail
		json.Unmarshal(j, &f)
		log.Println(f)
		os.Exit(2)
	}

	fmt.Printf("%#v\n", w)
}
