package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGpChapter(t *testing.T) {
	fs := strings.Split(
		"GOPR0017.MP4 GOPR0022.MP4 GP010018.MP4 GP020017.MP4 GP020022.MP4 GP030022.MP4 GOPR0018.MP4  GP010017.MP4  GP010022.MP4  GP020018.MP4  GP030017.MP4  GP040017.MP4",
		" ",
	)

	got := gpChapter(fs)

	expect := map[string][]string{
		"0022": strings.Split("GOPR0022.MP4 GP010022.MP4 GP020022.MP4 GP030022.MP4", " "),
		"0018": strings.Split("GOPR0018.MP4 GP010018.MP4 GP020018.MP4", " "),
		"0017": strings.Split("GOPR0017.MP4 GP010017.MP4 GP020017.MP4 GP030017.MP4 GP040017.MP4", " "),
	}

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("\nexpect %v\ngot %v\n", expect, got)
	}
}
