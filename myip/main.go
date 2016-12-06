package main

import (
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
	fPort              = flag.String("p", "8081", "server port")
	fServer            = flag.Bool("s", false, "print IP from client")
	fSpeedTestDuration = flag.Int("d", 0, "speed test duration in secs")
	fWaitBeforeExit    = flag.Bool("w", false, "wait userinput before exit")
)

func main() {
	flag.Parse()

	if *fServer { // server
		l, err := net.Listen("tcp", ":"+*fPort)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		log.Println("Listening on", *fPort)
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}

			log.Printf("Received message %s -> %s\n",
				conn.RemoteAddr(), conn.LocalAddr())

			go func(conn net.Conn) {
				defer conn.Close()

				// if *flagSpeedTestDuration == 0 {
				// 	_, err := io.Copy(os.Stdout, conn)
				// 	if err != nil {
				// 		log.Println("error at read:", err)
				// 	}
				// } else {
				log.Println("waiting for all received...")
				tr := tachoio.NewReader(conn)
				io.Copy(ioutil.Discard, tr)
				n, d := tr.ReadMeter()
				log.Printf("read %s bytes per a sec\n", scale(int(float64(n)/d.Seconds())))
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
			addr := serverIP + ":" + *fPort
			log.Println("Send ip to", addr)

			c, err = net.Dial("tcp", addr)
			if err != nil {
				panic(err)
			}
		} else {
			c = os.Stdout
		}
		defer c.Close()

		if *fSpeedTestDuration == 0 {
			fmt.Fprintf(c, "IP: %s\nMAC: %s\n", ip, mac)
		} else {
			secTick := time.Tick(time.Second)
			tw := tachoio.NewWriter(c)
			go func() { io.Copy(tw, &tachoio.NoopRead{}) }()
			for i := 0; i < *fSpeedTestDuration; i++ {
				select {
				case <-secTick:
					n, d := tw.WriteMeter()
					log.Printf("write %s bytes per a sec\n", scale(int(float64(n)/d.Seconds())))
				}
			}
		}
	}

	if *fWaitBeforeExit {
		var s string
		fmt.Scanln(&s)
	}
}
