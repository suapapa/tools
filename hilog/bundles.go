// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var bundles = map[string]logSchemePack{
	"gotest": logSchemePack{
		{
			Type: "re", Ptn: "^ok",
			FG: "GREEN",
		}, {
			Type: "re", Ptn: "^FAIL",
			BG: "RED",
		}, {
			Type: "startswith", Ptn: "--- PASS:",
			FG: "BLUE",
		}, {
			Type: "startswith", Ptn: "--- FAIL:",
			FG: "RED",
		},
	},
	"android-logcat": logSchemePack{
		{
			Type: "startswith", Ptn: "F/",
			FG: "YELLOW", BG: "RED",
			Attrs: []string{"BOLD"},
		}, {
			Type: "startswith", Ptn: "E/",
			FG:    "RED",
			Attrs: []string{"BOLD"},
		}, {
			Type: "startswith", Ptn: "W/",
			FG:    "YELLOW",
			Attrs: []string{"BOLD"},
		}, {
			Type: "startswith", Ptn: "I/",
			FG:    "WHITE",
			Attrs: []string{"BOLD"},
		}, {
			Type: "startswith", Ptn: "D/",
			FG: "GREEN",
		}, {
			Type: "startswith", Ptn: "V/",
			FG: "BLUE",
		},
	},
	"linux-kmsg": logSchemePack{
		{
			// emergency, alert, crit
			Type: "re", Ptn: "^<[012]>",
			FG: "YELLOW", BG: "RED",
			Attrs: []string{"BOLD"},
		}, {
			// error
			Type: "startswith", Ptn: "<3>",
			FG:    "RED",
			Attrs: []string{"BOLD"},
		}, {
			// warning
			Type: "startswith", Ptn: "<4>",
			FG:    "YELLOW",
			Attrs: []string{"BOLD"},
		}, {
			// notice
			Type: "startswith", Ptn: "<5>",
			FG:    "WHITE",
			Attrs: []string{"BOLD"},
		}, {
			// info
			Type: "startswith", Ptn: "<6>",
			FG: "GREEN",
		}, {
			// debug
			Type: "startswith", Ptn: "<7>",
			FG: "BLUE",
		},
	},
	"red-stderr": logSchemePack{
		{
			Source: "stderr",
			Type:   "re", Ptn: ".",
			FG: "RED",
		},
	},
}
