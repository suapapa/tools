// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
)

var (
	lsp  logSchemePack
	errC = make(chan error, 2)
)

func main() {
	opt, args := parseFlags()

	if opt.listBundledLSPs {
		// Print out bundled LSPs
		for b := range bundles {
			fmt.Println(b)
		}
		os.Exit(0)
	}

	var err error
	lsp, err = loadLogSchemePack(args[0])
	exitIf(err)

	if opt.dumpLSP {
		// Print out lsp in JSON format
		lsp.MarshalToJson(os.Stdout)
		os.Exit(0)
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	args = args[1:]
	if len(args) > 0 {
		// command given as args
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cStdOut, err := cmd.StdoutPipe()
		exitIf(err)
		cStdErr, err := cmd.StderrPipe()
		exitIf(err)

		err = cmd.Start()
		exitIf(err)

		go paintLines(ctx, cStdOut, "stdout")
		go paintLines(ctx, cStdErr, "stderr")

		err = <-errC
		exitIf(err)
		err = <-errC
		exitIf(err)

		cmd.Wait()
	} else {
		// log is given through stdin
		go paintLines(ctx, os.Stdin, "stdout")

		err = <-errC
		exitIf(err)
	}
}

func paintLines(ctx context.Context, r io.Reader, source string) {
	scanner := bufio.NewScanner(r)
	doneC := ctx.Done()
	for scanner.Scan() {
		lsp.Println(source, scanner.Text())
		select {
		case <-doneC:
			break
		default:
			// do nothing
		}
	}

	if err := ctx.Err(); err != nil {
		errC <- err
		return
	}

	if err := scanner.Err(); err != nil {
		errC <- err
		return
	}

	errC <- nil
}

func loadLogSchemePack(lspName string) (logSchemePack, error) {
	// Load from current directory
	if lsp := loadLogSchemePackFile(lspName + ".json"); lsp != nil {
		return lsp, nil
	}

	// Load from ~/configs/.hilog
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("[hilog] %v", err)
	}

	lspPath := path.Join(usr.HomeDir, "configs", ".hilog", lspName+".json")
	if lsp := loadLogSchemePackFile(lspPath); lsp != nil {
		return lsp, nil
	}

	// Load from bundles
	if lsp, ok := bundles[lspName]; ok {
		return lsp, nil
	}

	return nil, fmt.Errorf("[hilog] fail to load LogSchemePack, %s", lspName)
}

func exitIf(err error) {
	if err != nil {
		panic(err)
	}
}
