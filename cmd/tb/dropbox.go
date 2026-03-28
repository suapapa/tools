// +build !windows

package main

import (
	"os"
	"path/filepath"
)

func dropboxDir() string {
	tomboyRoot := filepath.Join(os.Getenv("HOME"), "Dropbox")
	return tomboyRoot
}
