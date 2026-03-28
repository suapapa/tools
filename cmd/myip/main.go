package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/suapapa/tools/pkg/scale"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *fWaitBeforeExit {
		fmt.Println("Press Enter to exit...")
		var s string
		fmt.Scanln(&s)
	}
}

func run(ctx context.Context) error {
	if *fServer {
		return runServer(ctx)
	}
	return runClient(ctx)
}

func runServer(ctx context.Context) error {
	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", ":"+*fPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer l.Close()

	log.Printf("Listening on :%s", *fPort)

	go func() {
		<-ctx.Done()
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Printf("failed to accept: %v", err)
				continue
			}
		}

		go handleServerConn(conn)
	}
}

func handleServerConn(conn net.Conn) {
	defer conn.Close()
	bc := bufio.NewReader(conn)
	firstline, err := bc.ReadString('\n')
	if err == nil {
		firstline = strings.TrimRight(firstline, "\r\n")
		fmt.Print(firstline)
	}

	tr := tachoio.NewReader(bc)
	_, _ = io.Copy(io.Discard, tr)
	n, d := tr.ReadMeter()
	
	bps := int(float64(n) / d.Seconds())
	formatted, _ := scale.FormatWithDecimal(strconv.Itoa(bps), scale.Bin)
	fmt.Printf(" %s BPS\n", formatted)
}

func runClient(ctx context.Context) error {
	ip, err := resolveIP()
	if err != nil {
		return fmt.Errorf("failed to resolve IP: %w", err)
	}

	var c io.WriteCloser
	serverIP := flag.Arg(0)
	if serverIP != "" {
		addr := serverIP + ":" + *fPort
		log.Printf("Sending IP to %s", addr)

		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to dial: %w", err)
		}
		c = conn
	} else {
		c = os.Stdout
	}
	defer c.Close()

	fmt.Fprintln(c, ip)

	if *fSpeedTestDuration > 0 {
		return runSpeedTest(ctx, c)
	}

	return nil
}

func runSpeedTest(ctx context.Context, w io.Writer) error {
	tw := tachoio.NewWriter(w)
	go func() {
		_, _ = io.Copy(tw, tachoio.NoopReader)
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	duration := time.Duration(*fSpeedTestDuration) * time.Second
	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		case <-ticker.C:
			n, d := tw.WriteMeter()
			bps := int(float64(n) / d.Seconds())
			formatted, _ := scale.FormatWithDecimal(strconv.Itoa(bps), scale.Bin)
			log.Printf("write %s bytes per a sec", formatted)
		}
	}
}
