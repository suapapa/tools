// +build windows

package main

import (
	"os/user"
	"path/filepath"
)

func dropboxDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, "Dropbox")
}
