package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
)

func main() {
	fPort := flag.String("p", "8080", "fPort")
	fDir := flag.String("d", "./", "directory")
	flag.Parse()

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, err := filepath.Abs(*fDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("serving", dir, "...")
	// fmt.Println(dir)

	ipStr, err := resolveIP()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("http://" + ipStr + ":" + *fPort)

	log.Fatal(http.ListenAndServe(":"+*fPort, http.FileServer(http.Dir(dir))))
}

func resolveIP() (string, error) {
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	return "", err
	// }

	var ip net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

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
			if ip != nil {
				ip = ip.To4()
				return ip.String(), nil
				// return strings.Join([]string{hostname, ip.String(), i.HardwareAddr.String()}, " "), nil
			}
		}
	}
	return "", fmt.Errorf("cannot resolve the IP")
}
