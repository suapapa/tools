package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	iconv "github.com/djimenez/iconv-go"
)

var encFrom string
var encTo string
var dryRun bool
var converter *iconv.Converter

func init() {
	flag.StringVar(&encFrom, "f", "auto", "encording of input zip file")
	flag.StringVar(&encTo, "t", "auto", "encording of the environment")
	flag.BoolVar(&dryRun, "n", false, "dry-run")
	flag.Parse()

	var sysLocale, sysEncording string

	envLANG := os.Getenv("LANG")
	if envLANG == "" {
		sysLocale = "unknown"
		sysEncording = "unknown"
	} else {
		s := strings.SplitN(envLANG, ".", 2) // example ko_KR.utf8
		sysLocale = s[0]
		sysEncording = s[1]
	}

	if encFrom == "auto" {
		// guess native encording by system locale
		// TODO: add more native encordings for locales
		switch sysLocale {
		case "ko_KR":
			encFrom = "cp949"
		default:
			encFrom = "ascii"
		}
	}

	if encTo == "auto" {
		switch sysEncording {
		case "unknown":
			encTo = "utf8"
		default:
			encTo = sysEncording
		}
	}

	fmt.Println("conv from", encFrom, "to", encTo)
	converter, _ = iconv.NewConverter(encFrom, encTo)
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
		fixed, _ := converter.ConvertString(f.Name)
		fmt.Println(fixed)
		if dryRun {
			continue
		}

		err := os.MkdirAll(filepath.Dir(fixed), os.ModePerm)
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
