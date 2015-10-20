// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"
)

func TestConvertForHuman(t *testing.T) {
	numbers := []struct {
		first string
		zeros int

		expect string
	}{
		{"2", 2, "200"},
		{"3", 4, "3만"},
		{"3", 6, "300만"},
		{"3", 10, "300억"},
		{"3", 14, "300조"},
		{"3", 18, "300경"},
		{"3", 22, "300해"},

		{"3", 3, "3000"},
		{"3", 7, "3000만"},
		{"3", 11, "3000억"}, // len(in) == 12
	}

	for _, n := range numbers {
		f, zs, exp := n.first, n.zeros, n.expect
		in := f + strings.Repeat("0", zs)

		if r := convertForHuman(in, &koScale); r != exp {
			t.Errorf("want \"%s\" got \"%s\"", exp, r)
		}
	}
}
