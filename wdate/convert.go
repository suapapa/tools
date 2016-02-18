// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strconv"

// +0700 -> 07*60*60 + 30*60
func timeOffsetStrToInt(s string) (int, error) {
	if len(s) != 5 {
		return 0, errInvalidOffset
	}

	if s[0] != '+' && s[0] != '-' {
		return 0, errInvalidOffset
	}

	hStr, mStr := s[1:3], s[3:5]
	h, err := strconv.Atoi(hStr)
	if err != nil {
		return 0, err
	}

	m, err := strconv.Atoi(mStr)
	if err != nil {
		return 0, err
	}
	switch s[0] {
	case '+':
		return h*60*60 + m*60, nil
	case '-':
		return -h*60*60 - m*60, nil
	}
	return 0, errInvalidOffset
}
