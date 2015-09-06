// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

type testeePack []struct {
	str    string
	expect bool
}

func runTests(t *testing.T, ls *logScheme, testees testeePack) {
	for _, testee := range testees {
		if ls.Match("", testee.str) != testee.expect {
			t.Errorf("failed to match %s on %v",
				ls.Type, testee)
		}
	}
}

func TestMatchRe(t *testing.T) {
	rePtn := logScheme{
		Type: "re",
		Ptn:  ".*ffle$",
	}

	testees := testeePack{
		{"Belgian Waffle", true},
		{"Shuffle", true},
		{"Shuffle dance", false},
	}

	runTests(t, &rePtn, testees)
}

func TestMatchStartswith(t *testing.T) {
	startsWith := logScheme{
		Type: "startswith",
		Ptn:  "cat",
	}

	testees := testeePack{
		{"cat tower", true},
		{"bat", false},
		{"logcat", false},
	}

	runTests(t, &startsWith, testees)
}

func TestMatchEndswith(t *testing.T) {
	endsWith := logScheme{
		Type: "endswith",
		Ptn:  "jack",
	}

	testees := testeePack{
		{"apple jack", true},
		{"jacky chan", false},
		{"bluce lee", false},
	}

	runTests(t, &endsWith, testees)
}

func TestMatchContains(t *testing.T) {
	contains := logScheme{
		Type: "contains",
		Ptn:  "bene",
	}

	testees := testeePack{
		{"cafe bene", true},
		{"beneton", true},
		{"bean pole", false},
	}

	runTests(t, &contains, testees)
}
