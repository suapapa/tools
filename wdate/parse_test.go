// Copyright 2017, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestParseDate(t *testing.T) {
	// TODO: add test
	testCases := []struct {
		dateStr string
		yearDay int
		weekCnt int
	}{
		{"2017-01-01", 1, 1}, // Weired... I thinks its 0th week
		{"2017-01-02", 2, 1}, {"2017-01-08", 8, 1},
		{"2017-01-09", 9, 2}, {"2017-01-15", 15, 2},
	}

	for _, tc := range testCases {
		tt, err := parseTimeStr(tc.dateStr)
		if err != nil {
			panic(err)
		}

		// // test YearDay
		// if tc.yearDay != tt.YearDay() {
		// 	t.Errorf("yd wrong! exp:%d, got:%d\n",
		// 		tc.yearDay, tt.YearDay())
		// }

		// test weekCount
		// log.Println("test", tt, weekCount(tt))
		if tc.weekCnt != weekCount(tt) {
			t.Errorf("wc wrong! exp:%d, got:%d %v\n",
				tc.weekCnt, weekCount(tt), tt)
		}
	}
}
