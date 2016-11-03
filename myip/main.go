package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	if flagServer { // server
		l, err := net.Listen("tcp", ":"+flagPort)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		log.Println("Listening on", flagPort)
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}

			log.Printf("Received message %s -> %s\n",
				conn.RemoteAddr(), conn.LocalAddr())

			go handleRequest(conn)
		}
	} else { // client
		ip, err := resolveIP()
		if err != nil {
			panic(err)
		}

		var c io.WriteCloser

		serverIP := flag.Arg(0)
		if serverIP != "" {
			addr := serverIP + ":" + flagPort
			log.Println("Send ip to", addr)

			c, err = net.Dial("tcp", addr)
			if err != nil {
				panic(err)
			}
		} else {
			c = os.Stdout
		}

		fmt.Fprintf(c, "IP: %s", ip)
		c.Close()
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.Println("error at read:", err)
		return
	}

	fmt.Println(string(buf[:reqLen]))
}
