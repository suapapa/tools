// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var (
	koScale = Scale{
		n: 4,
		name: []string{
			"",
			"만", // 10 ** 4
			"억", // 10 ** 8
			"조", // 10 ** 12
			"경", // 10 ** 16
			"해", // 10 ** 20
		},
	}

	enScale = Scale{
		n: 3,
		name: []string{
			"",
			"Thousand",
			"Million",
			"Billion",
			"Quadrillion",
			"Quintillion",
			"Sextillion",
			"Septillion",
		},
	}

	binScale = Scale{
		n: 3,
		name: []string{
			"",
			"Kilo",
			"Mega",
			"Giga",
			"Tera",
			"Peta",
			"Exa",
			"Zetta",
			"yotta",
		},
	}
)
