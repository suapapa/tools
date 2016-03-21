package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/suapapa/go_hangul/encoding/cp949"
)

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "n", false, "dry-run")
	flag.Parse()
}

func checkIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func convUnzip(fileName string) {
	r, err := zip.OpenReader(fileName)
	checkIf(err)
	defer r.Close()

	for _, f := range r.File {
		fixedBytes, err := cp949.From([]byte(f.Name))
		checkIf(err)
		fixed := string(fixedBytes)
		fmt.Println(fixed)
		if dryRun {
			continue
		}

		err = os.MkdirAll(filepath.Dir(fixed), os.ModePerm)
		checkIf(err)

		if f.FileInfo().IsDir() {
			continue
		}

		o, err := os.Create(fixed)
		checkIf(err)

		r, err := f.Open()
		checkIf(err)

		io.Copy(o, r)
		r.Close()
		o.Close()
	}
}

func main() {
	for _, zipFilename := range flag.Args() {
		fmt.Println("unzipping", zipFilename, "...")
		convUnzip(zipFilename)
	}
}
