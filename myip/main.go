package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	tachoio "github.com/suapapa/go_tachoio"
)

var (
	flagPort              string
	flagServer            bool
	flagSpeedTestDuration int
)

func init() {
	flag.StringVar(&flagPort, "p", "8081", "server port")
	flag.BoolVar(&flagServer, "s", false, "print IP from client")
	flag.IntVar(&flagSpeedTestDuration, "d", 0, "speed test duration in secs")
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

			go func(conn net.Conn) {
				defer conn.Close()

				// if flagSpeedTestDuration == 0 {
				// 	_, err := io.Copy(os.Stdout, conn)
				// 	if err != nil {
				// 		log.Println("error at read:", err)
				// 	}
				// } else {
				log.Println("waiting for all received...")
				tr := tachoio.NewReader(conn)
				io.Copy(ioutil.Discard, tr)
				n, d := tr.ReadMeter()
				log.Printf("read %.1f bytes per a sec\n", float64(n)/d.Seconds())
				// }
			}(conn)
		}
	} else { // client
		ip, mac, err := resolveIP()
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
		defer c.Close()

		if flagSpeedTestDuration == 0 {
			fmt.Fprintf(c, "IP: %s\nMAC: %s\n", ip, mac)
		} else {
			secTick := time.Tick(time.Second)
			tw := tachoio.NewWriter(c)
			go func() { io.Copy(tw, rand.Reader) }()
			for i := 0; i < flagSpeedTestDuration; i++ {
				select {
				case <-secTick:
					n, d := tw.WriteMeter()
					log.Printf("write %.1f bytes per a sec\n", float64(n)/d.Seconds())
				}
			}
		}
	}
}
