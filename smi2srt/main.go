// Copyright 2015, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/saintfish/chardet"
	"github.com/suapapa/go_hangul/encoding/cp949"
	"github.com/suapapa/go_subtitle"
)

var (
	flagOut string
)

func init() {
	flag.StringVar(&flagOut, "o", "", "srt file name")
}

func main() {
	flag.Parse()
	var err error

	// read smi
	var s *os.File
	if flag.NArg() == 0 {
		s = os.Stdin
	} else {
		smiFileName := flag.Arg(0)
		s, err = os.Open(smiFileName)
		if err != nil {
			panic(err)
		}
		defer s.Close()
	}

	// detect charset
	buff, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}
	det := chardet.NewHtmlDetector()
	r, err := det.DetectBest(buff)
	var utf8Buff []byte
	switch r.Charset {
	case "EUC-KR":
		utf8Buff, err = cp949.From(buff)
	case "UTF-8":
		utf8Buff = buff
	default:
		log.Fatal(r.Charset, "is not supported")
	}

	utf8Src := bytes.NewBuffer(utf8Buff)

	// convert it to subtitle.Book
	b, err := subtitle.ReadSmi(utf8Src)
	if err != nil {
		panic(err)
	}

	// write in srt
	var w *os.File
	if flagOut != "" {
		w, err = os.Create(flagOut)
		if err != nil {
			panic(err)
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	err = subtitle.ExportToSrtFile(b, w)
	if err != nil {
		panic(err)
	}
}
