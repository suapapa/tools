package main

import (
	"os/exec"
	"strings"
)

func runAdbDevices() []string {
	stdout, _ := runCommandSync("adb devices")

	lines := strings.Split(stdout, "\n")
	ret := make([]string, 0)
	for _, l := range lines {
		if strings.HasPrefix(l, "List of devices attached") {
			continue
		}
		if strings.HasPrefix(l, "*") {
			continue
		}
		if strings.TrimSpace(l) == "" {
			continue
		}

		// TODO: filter unauthorized?
		d := strings.Split(l, "\t")
		ret = append(ret, d[0])
	}
	return ret
}

func runCommandSync(cmdStr string) (string, error) {
	cmdSlice := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}
