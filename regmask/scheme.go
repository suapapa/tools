// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
)

func parseCsvFile(fileName string) (bitsFields, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	p := make(bitsFields, 0, len(records))

	// csv is formed;
	// lable, offset, cnt
	for _, r := range records {
		if len(r) != 3 {
			return nil, errors.New("invalid format")
		}
		offset, err := strconv.Atoi(r[1])
		if err != nil {
			return nil, err
		}
		cnt, err := strconv.Atoi(r[2])
		if err != nil {
			return nil, err
		}

		f, err := newBitsField(r[0], offset, cnt)
		if err != nil {
			return nil, err
		}

		p = append(p, f)
	}

	return p, nil
}
