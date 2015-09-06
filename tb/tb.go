package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/suapapa/tools/tb/tomboy"
)

var (
	flagID string
)

func main() {
	flag.StringVar(&flagID, "i", "", "show a note of given ID")
	flag.Parse()

	if flagID == "" && len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	var tomboyRoot = filepath.Join(os.Getenv("HOME"), "/Dropbox/tomboy/")
	notes, err := tomboy.MakeNotebookFromFileSystemSync(tomboyRoot)
	if err != nil {
		panic(err)
	}

	// display
	if flagID != "" {
		for _, n := range notes {
			if strings.HasPrefix(n.ID.String(), flagID) {
				fmt.Println(n.Text.Content)
			}
		}
		os.Exit(0)
	}

	// search
	searchKey := flag.Arg(0)
	if searchKey == "" {
		os.Exit(0)
	}

	var searchedNotes tomboy.Notebook
	for _, n := range notes {
		if strings.Contains(n.Title, searchKey) {
			searchedNotes = append(searchedNotes, n)
			continue
		}
		if strings.Contains(n.Text.Content, searchKey) {
			searchedNotes = append(searchedNotes, n)
		}
	}

	for _, n := range searchedNotes {
		fmt.Println(n.ID, n.Title, n.Rev)
	}

}
