// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
)

type bitsField struct {
	lable  string
	offset int
	cnt    int
	mask   uint64
}

func newBitsField(lable string, offset, cnt int) (*bitsField, error) {
	if offset < 0 || cnt < 0 {
		return nil, errors.New("offset and cnt must be posive")
	}

	f := &bitsField{
		lable:  lable,
		offset: offset,
		cnt:    cnt,
	}
	for i := 0; i < cnt; i++ {
		f.mask |= 1 << uint32(i+offset)
	}
	return f, nil
}

// Value returns bits-field of value given value
func (f *bitsField) Value(v uint64) uint64 {
	return (v & f.mask) >> uint32(f.offset)
}

func (f *bitsField) String() string {
	if opt.headerGuide {
		if f.cnt == 1 {
			return fmt.Sprintf("%s[%d]", f.lable, f.offset)
		}
		return fmt.Sprintf("%s[%d:%d]", f.lable, f.cnt+f.offset-1, f.offset)
	}
	return f.lable
}

type bitsFields []*bitsField

func (p bitsFields) Strings() (s []string) {
	for _, f := range p {
		s = append(s, f.String())
	}
	return
}

func (p bitsFields) String() (s string) {
	for i, f := range p {
		if i != 0 {
			s += ","
		}
		s += f.String()
	}
	return
}
