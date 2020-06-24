package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	port    int
	count   int
	serveFN string
	quitCh  = make(chan struct{})
)

func init() {
	flag.IntVar(&port, "p", 8080, "listen port")
	flag.IntVar(&count, "c", 1, "how many download")
}

func main() {
	flag.Parse()
	serveFN = flag.Arg(0)
	if serveFN == "" {
		log.Println("should give a file to server")
		os.Exit(-1)
	}

	ip, err := resolveIP()
	if err != nil {
		panic(err)
	}
	log.Printf("listen to http://%s:%d/%s", ip, port, filepath.Base(serveFN))

	srv := startHTTPServer()
	<-quitCh
	err = srv.Shutdown(context.TODO())
	if err != nil {
		panic(err)
	}
	log.Println("bye~")
}

func startHTTPServer() *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("access from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(serveFN)))
		f, err := os.Open(serveFN)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		io.Copy(w, f)
		count--
		if count <= 0 {
			quitCh <- struct{}{}
		}
	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	return srv
}

// resolveIP returns hostname, IP, MAC and error
func resolveIP() (string, error) {
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
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("cannot resolve the IP")
}
