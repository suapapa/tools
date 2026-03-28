package main

import "testing"

func TestConvertExtToSrt(t *testing.T) {
	in := "script.smi"
	expect := "script.srt"
	got := convertExtToSrt(in)

	if got != expect {
		t.Errorf("expect %s but got %s", expect, got)
	}
}
