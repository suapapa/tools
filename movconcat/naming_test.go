package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGpChapter(t *testing.T) {
	fs := strings.Split("GOPR0021.MP4 GP010021.MP4 GP020021.MP4 GP030021.MP4 GP040021.MP4", " ")

	got := gpChapter(fs)

	expect := map[string][]string{
		"0021": []string{"GOPR0021.MP4", "GP010021.MP4", "GP020021.MP4", "GP030021.MP4", "GP040021.MP4"},
	}

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("expect %v, got %v\n", expect, got)
	}
}
