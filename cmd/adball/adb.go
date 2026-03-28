package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func runAdbDevices(ctx context.Context) ([]string, error) {
	stdout, err := runCommand(ctx, "adb", "devices")
	if err != nil {
		return nil, fmt.Errorf("failed to run adb devices: %w", err)
	}

	lines := strings.Split(stdout, "\n")
	ret := make([]string, 0)
	for _, l := range lines {
		if strings.HasPrefix(l, "List of devices attached") {
			continue
		}
		if strings.HasPrefix(l, "*") {
			continue
		}
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		// TODO: filter unauthorized?
		d := strings.Split(l, "\t")
		if len(d) > 0 {
			ret = append(ret, d[0])
		}
	}
	return ret, nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
