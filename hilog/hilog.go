// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	opt, args := parseFlags()

	if opt.listBundledLSPs {
		// Print out bundled LSPs
		for lsp := range bundles {
			fmt.Println(lsp)
		}
		os.Exit(0)
	}

	lsp, err := loadLogSchemePack(args[0])
	exitIf(err)

	if opt.dumpLSP {
		// Print out lsp in JSON format
		lsp.MarshalToJson(os.Stdout)
		os.Exit(0)
	}

	args = args[1:]
	if len(args) > 0 {
		// command given as args
		cmd := exec.Command(args[0], args[1:]...)
		cStdOut, err := cmd.StdoutPipe()
		exitIf(err)
		cStdErr, err := cmd.StderrPipe()
		exitIf(err)

		wg.Add(2)
		go paintLines(cStdOut, lsp, "stdout")
		go paintLines(cStdErr, lsp, "stderr")

		cmd.Start()
		wg.Wait()
		cmd.Wait()
	} else {
		// output given as stdin
		wg.Add(1)
		paintLines(os.Stdin, lsp, "stdout")
		wg.Wait()
	}

}

func paintLines(r io.Reader, lsp logSchemePack, source string) {
	defer wg.Done()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lsp.Println(source, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[hilog] err of paintLines on", source, err)
	}
}

func loadLogSchemePack(lspName string) (logSchemePack, error) {
	// Load from current directory
	if lsp := loadLogSchemePackFile(lspName + ".json"); lsp != nil {
		return lsp, nil
	}

	// Load from ~/configs/.hilog
	usr, err := user.Current()
	exitIf(err)

	lspPath := path.Join(usr.HomeDir, "configs", ".hilog", lspName+".json")
	if lsp := loadLogSchemePackFile(lspPath); lsp != nil {
		return lsp, nil
	}

	// Load from bundles
	if lsp, ok := bundles[lspName]; ok {
		return lsp, nil
	}

	return nil, fmt.Errorf("[hilog] Can not load LogSchemePack, %s.", lspName)
}

func exitIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
