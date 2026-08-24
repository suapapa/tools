// Copyright 2015, Homin Lee <i@homin.dev>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/suapapa/tools/pkg/scale"
)

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

			if r, err := scale.Format(in, scale.Ko); err != nil || r != exp {
				if err != nil {
					t.Fatalf("got err: %v", err)
				}
				t.Errorf("want %q got %q", exp, r)
			}
		})
	}
}

