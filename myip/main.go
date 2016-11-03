package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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
	} else {
		ip, err := resolveIP()
		if err != nil {
			panic(err)
		}

		serverIP := flag.Arg(0)

		if serverIP != "" {
			fmt.Println("serverIP", serverIP)
			// TODO: Send ip to server
			// echo "hello" | nc 127.0.0.1 8081
		}

		fmt.Println("IP:", ip)
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
