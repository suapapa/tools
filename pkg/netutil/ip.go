package netutil

import (
	"fmt"
	"net"
)

// IPInfo contains network interface information.
type IPInfo struct {
	Interface string
	IP        string
	MAC       string
}

// GetExternalIPs returns information about non-loopback, up interfaces.
func GetExternalIPs() ([]IPInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	var infos []IPInfo
	for _, i := range ifaces {
		if (i.Flags&net.FlagUp) == 0 || (i.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			ip = ip.To4()
			if ip != nil {
				infos = append(infos, IPInfo{
					Interface: i.Name,
					IP:        ip.String(),
					MAC:       i.HardwareAddr.String(),
				})
			}
		}
	}

	if len(infos) == 0 {
		return nil, fmt.Errorf("no external IP found")
	}

	return infos, nil
}

// FirstExternalIP returns the first non-loopback IPv4 address found.
func FirstExternalIP() (string, error) {
	ips, err := GetExternalIPs()
	if err != nil {
		return "", err
	}
	return ips[0].IP, nil
}
