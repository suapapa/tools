package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/suapapa/tools/tb/tomboy"
)

var (
	flagListAll    bool
	flagID         string
	flagWithMarkup bool
)

var reMakrup = regexp.MustCompile("</?[^<>]+?>")

func main() {
	flag.BoolVar(&flagListAll, "a", false, "list all note")
	flag.StringVar(&flagID, "i", "", "show a note of given ID")
	flag.BoolVar(&flagWithMarkup, "m", false, "show a note with tomboy markup")
	flag.Parse()

	if flagID == "" && flagListAll == false && flag.NArg() == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	var tomboyRoot = filepath.Join(os.Getenv("HOME"), "/Dropbox/tomboy/")
	notes, err := tomboy.MakeNotebookFromFileSystemSync(tomboyRoot)
	if err != nil {
		panic(err)
	}

	// display a note
	if flagID != "" {
		for _, n := range notes {
			if strings.HasPrefix(n.ID.String(), flagID) {
				if flagWithMarkup {
					fmt.Println(n.Text.Content)
				} else {
					noteText := reMakrup.ReplaceAllString(n.Text.Content, "")
					fmt.Println(noteText)
				}
			}
		}
		os.Exit(0)
	}

	// list All
	if flagListAll {
		for _, n := range notes {
			fmt.Println(n.ID, n.Title, n.Rev)
		}

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
