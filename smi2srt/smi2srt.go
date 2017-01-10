package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/saintfish/chardet"
	"github.com/suapapa/go_hangul/encoding/cp949"
	subtitle "github.com/suapapa/go_subtitle"
)

// read smi from r and write srt to w
func smi2srt(r io.Reader, w io.Writer) error {
	// detect charset
	buff, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	det := chardet.NewHtmlDetector()
	d, err := det.DetectBest(buff)
	if err != nil {
		return err
	}

	// cp949 to utf8
	var utf8Buff []byte
	switch d.Charset {
	case "EUC-KR":
		utf8Buff, err = cp949.From(buff)
	case "UTF-8":
		utf8Buff = buff
	default:
		return fmt.Errorf("smi2srt: encoding, %s is not supported", d.Charset)
	}

	utf8Src := bytes.NewBuffer(utf8Buff)

	// convert it to subtitle.Book
	b, err := subtitle.ReadSmi(utf8Src)
	if err != nil {
		return err
	}

	// TODO: fix subtitle package to srt.Encode(...)
	err = subtitle.ExportToSrtFile(b, w)
	if err != nil {
		panic(err)
	}

	return nil
}
