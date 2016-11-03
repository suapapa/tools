package main

import (
	"fmt"
	"net"
)

func resolveIP() (string, error) {
	var ip net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ip = ip.To4()
			if ip != nil && ip.String() != "127.0.0.1" {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("cannot resolve the IP")
}
