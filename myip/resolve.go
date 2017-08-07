package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// resolveIP returns hostname, IP, MAC and error
func resolveIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var ip net.IP
	for _, i := range ifaces {
		if (i.Flags&net.FlagUp) == 0 ||
			(i.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// sometimes ip.To4() makes ip to nil
			if ip != nil {
				ip = ip.To4()
			}
			if ip != nil {
				return strings.Join([]string{hostname, ip.String(), i.HardwareAddr.String()}, " "), nil
			}
		}
	}
	return "", fmt.Errorf("cannot resolve the IP")
}
