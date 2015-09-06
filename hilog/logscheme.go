// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

type logSchemePack []logScheme

func (lsp logSchemePack) MarshalToJson(w io.Writer) {
	b, err := json.MarshalIndent(lsp, "", "    ")
	if err != nil {
		return
	}
	w.Write(b)
}

func (lsp logSchemePack) Println(source, line string) {
	line = strings.TrimRight(line, "\n")
	for _, ls := range lsp {
		if ls.Match(source, line) {
			ls.Println(line)
			return
		}
	}
	fmt.Println(line)
}

func loadLogSchemePackFile(filename string) (lsp logSchemePack) {
	if d, err := ioutil.ReadFile(filename); err == nil {
		json.Unmarshal(d, &lsp)
	}
	return
}

var colorMap = map[string]uint8{
	// "BLACK":   0,
	"RED":     1,
	"GREEN":   2,
	"YELLOW":  3,
	"BLUE":    4,
	"MAGENTA": 5,
	"CYAN":    6,
	"WHITE":   7,
	"BALCK":   8,
}

var attrMap = map[string]string{
	"BOLD":      "\x1b[1m",
	"UNDERLINE": "\x1b[4m",
	"BLINK":     "\x1b[5m",
	"INVERSE":   "\x1b[7m",
}

type logScheme struct {
	Ptn    string         // search pattern
	Type   string         // startswith | endswith | contains | re
	cRe    *regexp.Regexp // comiled regexp for Type re
	Source string         // stdin | stderr
	FG, BG string         // BALCK | REG | GREEN | YELLOW | BLUE | MAGENTA | CYAN | WHITE
	Attrs  []string       // BOLD | UNDERLINE | BLINK | INVERSE
}

func (s *logScheme) Println(str string) {
	var fgStr, bgStr, attrStr string

	fg := colorMap[s.FG]
	if fg != 0 {
		fgStr = fmt.Sprintf("\x1b[3%dm", fg)
	}

	bg := colorMap[s.BG]
	if bg != 0 {
		bgStr = fmt.Sprintf("\x1b[4%dm", bg)
	}

	for _, strAttr := range s.Attrs {
		attrStr += attrMap[strAttr]
	}

	fmt.Printf("%s%s%s%s\x1b[0m\n", fgStr, bgStr, attrStr, str)
}

func (s *logScheme) Match(source, str string) bool {
	if s.Source != "" && source != s.Source {
		return false
	}

	switch s.Type {
	case "startswith":
		return strings.HasPrefix(str, s.Ptn)
	case "endswith":
		return strings.HasSuffix(str, s.Ptn)
	case "contains":
		return strings.Contains(str, s.Ptn)
	case "re":
		// compile regexp once
		if s.cRe == nil {
			c, err := regexp.Compile(s.Ptn)
			if err != nil {
				return false
			}
			s.cRe = c
		}
		return s.cRe.MatchString(str)
	}

	// shouldn't arrive here!
	return false
}
