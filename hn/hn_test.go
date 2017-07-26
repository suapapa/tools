// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestKoScale(t *testing.T) {
	numbers := []struct {
		in     string
		expect string
	}{
		{"200", "200"},
		{"30000", "3만"},
		{"3000000", "300만"},
		{"30000000000", "300억"},
		{"300000000000000", "300조"},
		{"3000000000000000000", "300경"},
		{"30000000000000000000000", "300해"},

		{"3000", "3000"},
		{"30000000", "3000만"},
		{"300000000000", "3000억"}, // len(in) == 12

		{"10**10", "100억"},
	}

	for _, n := range numbers {
		t.Run(n.expect, func(t *testing.T) {
			in, exp := n.in, n.expect

			if r, err := convertForHuman(in, &koScale); err != nil || r != exp {
				if err != nil {
					t.Fatal("got err:", err)
				}
				t.Errorf("want \"%s\" got \"%s\"", exp, r)
			}
		})
	}
}

func TestIsNum(t *testing.T) {
	tt := []struct {
		n string
		b bool
	}{
		{"10", true},
		{"a", false},
	}

	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			if isNum(tc.n) != tc.b {
				t.Errorf("exp %v but not", tc.b)
			}
		})
	}
}

func TestConvertPowerForm(t *testing.T) {
	tt := []struct {
		p   string
		exp string
	}{
		{"1**0", "1"},
		{"10**0", "1"},
		{"10**10", "10000000000"},
		{"a**b", "a**b"},
	}
	for _, tc := range tt {
		t.Run(tc.p, func(t *testing.T) {
			if got := convertPowerForm(tc.p); got != tc.exp {
				t.Errorf("exp %v but, got %v", tc.exp, got)
			}
		})
	}
}
