// Copyright -2015, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iomon

import (
	"io"
	"sync"
	"time"
)

// WriteMon is struct for monitoring io.Writer
type WriteMon struct {
	sync.Mutex
	cnt           uint64
	lastCheckTime time.Time
	w             io.Writer
}

// NewWriteMon returns *WriteMon which is initilaized with current time
func NewWriteMon(w io.Writer) *WriteMon {
	var m WriteMon
	m.lastCheckTime = time.Now()
	m.w = w
	return &m
}

// Check returns written byte count since last check
func (m *WriteMon) Check() (n uint64, lastTime time.Time) {
	m.Lock()
	defer m.Unlock()
	n, lastTime = m.cnt, m.lastCheckTime
	m.cnt = 0
	m.lastCheckTime = time.Now()
	return
}

func (m *WriteMon) Write(b []byte) (int, error) {
	n, err := m.w.Write(b)
	if err != nil {
		return n, err
	}
	m.Lock()
	defer m.Unlock()
	m.cnt += uint64(n)
	return n, nil
}

// ReadMon is struct for monitoring io.Readr
type ReadMon struct {
	sync.Mutex
	cnt           uint64
	lastCheckTime time.Time
	r             io.Reader
}

// NewReadMon returns *ReadMon which is initilaized with current time
func NewReadMon(r io.Reader) *ReadMon {
	var m ReadMon
	m.lastCheckTime = time.Now()
	m.r = r
	return &m
}

// Check returns read byte count since last check time and duration
func (m *ReadMon) Check() (n uint64, lastTime time.Time) {
	m.Lock()
	defer m.Unlock()
	n, lastTime = m.cnt, m.lastCheckTime
	m.cnt = 0
	m.lastCheckTime = time.Now()
	return
}

func (m *ReadMon) Read(b []byte) (int, error) {
	n, err := m.r.Read(b)
	if err != nil {
		return n, err
	}
	m.Lock()
	defer m.Unlock()
	m.cnt += uint64(n)
	return n, nil
}
