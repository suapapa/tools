package main

import (
	"fmt"
	"os"

	"github.com/suapapa/tools/pkg/netutil"
)

// resolveIP returns hostname, IP, MAC and error
func resolveIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	ips, err := netutil.GetExternalIPs()
	if err != nil {
		return "", err
	}

	// Just return the first one as before
	info := ips[0]
	return fmt.Sprintf("%s %s %s", hostname, info.IP, info.MAC), nil
}
