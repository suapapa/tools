package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		printUsage()
	}

	title := os.Args[1]
	msg := os.Args[2]
	userKey := os.Getenv("PUSHOVER_USERKEY")
	appToken := os.Getenv("PUSHOVER_APPTOKEN")

	formValues := url.Values{
		"token":   []string{appToken},
		"user":    []string{userKey},
		"title":   []string{title},
		"message": []string{msg},
	}

	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", formValues)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Printf("\t%s title message\n", os.Args[0])
	fmt.Printf("\tor %s message\n", os.Args[1])
	os.Exit(-1)
}
