package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
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

		// log.Println("Listening on", *fPort)
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}

			// log.Printf("Received message %s -> %s\n",
			// 	conn.RemoteAddr(), conn.LocalAddr())

			go func(conn net.Conn) {
				defer conn.Close()
				bc := bufio.NewReader(conn)
				firstline, err := bc.ReadString('\n')
				if err == nil {
					firstline = strings.TrimRight(firstline, "\n")
					fmt.Print(firstline)
				}

				tr := tachoio.NewReader(bc)
				io.Copy(ioutil.Discard, tr)
				n, d := tr.ReadMeter()
				fmt.Printf(" %s BPS\n", scale(int(float64(n)/d.Seconds())))
			}(conn)
		}
	} else { // client
		ip, err := resolveIP()
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

		fmt.Fprintln(c, ip)
		if *fSpeedTestDuration != 0 {
			secTick := time.Tick(time.Second)
			tw := tachoio.NewWriter(c)
			go func() { io.Copy(tw, tachoio.NoopReader) }()
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
