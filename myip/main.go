package main

import (
	"flag"
	"fmt"
	"net"
)

const (
	defaultPort = "8081"
)

var (
	flagPort   string
	flagServer bool
)

func init() {
	flag.StringVar(&flagPort, "p", "8081", "server port")
	flag.BoolVar(&flagServer, "s", false, "print client IP")
	flag.Parse()
}

func main() {
	if flagServer {
		// TODO: receive ip from client and print it
		panic("not implemented")

	} else {
		ip, err := resolveIP()
		if err != nil {
			panic(err)
		}

		serverIP := flag.Arg(0)

		if serverIP != "" {
			fmt.Println("serverIP", serverIP)
			// TODO: Send ip to server
		}

		fmt.Println("IP:", ip)
	}
}

func resolveIP() (string, error) {
	var ip net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
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
